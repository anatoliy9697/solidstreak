package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"
)

type pgRepo struct {
	ctx      context.Context
	pool     *pgxpool.Pool
	subPlans map[string]*subPkg.Plan
}

func initPGRepo(c context.Context, p *pgxpool.Pool, subPlans map[string]*subPkg.Plan) *pgRepo {
	return &pgRepo{c, p, subPlans}
}

func (r pgRepo) Create(s *subPkg.Subscription) error {
	sql := `
		INSERT INTO user_subscriptions (
			active,
			plan_code,
			start_date,
			finish_date,
			user_id,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := r.pool.QueryRow(
		r.ctx,
		sql,
		s.Active,
		s.PlanCode,
		s.StartDate,
		s.FinishDate,
		s.UserID,
		s.CreatedAt,
		s.UpdatedAt,
	).Scan(&s.ID)

	return err
}

func (r pgRepo) Update(s *subPkg.Subscription) error {
	sql := `
		UPDATE user_subscriptions SET
			active = $1,
			plan_code = $2,
			start_date = $3,
			finish_date = $4,
			user_id = $5,
			created_at = $6,
			updated_at = $7
		WHERE id = $8
	`
	_, err := r.pool.Exec(
		r.ctx,
		sql,
		s.Active,
		s.PlanCode,
		s.StartDate,
		s.FinishDate,
		s.UserID,
		s.CreatedAt,
		s.UpdatedAt,
		s.ID,
	)

	return err
}

func (r pgRepo) GetActiveByUserID(userID int64) (*subPkg.Subscription, error) {
	s := &subPkg.Subscription{}

	sql := `
		SELECT id, active, plan_code, start_date, finish_date, user_id, created_at, updated_at
		FROM user_subscriptions 
		WHERE 
			active IS TRUE
			AND user_id = $1
			AND start_date <= (NOW() AT TIME ZONE 'UTC')::date
			AND (finish_date IS NULL OR finish_date + interval '1 day' > (NOW() AT TIME ZONE 'UTC')::date)
		LIMIT 1
	`
	err := r.pool.QueryRow(
		r.ctx,
		sql,
		userID,
	).Scan(
		&s.ID,
		&s.Active,
		&s.PlanCode,
		&s.StartDate,
		&s.FinishDate,
		&s.UserID,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.NewNotFoundErr(fmt.Sprintf("couldn't find user subscription for user ID %v", userID))
		}
		return nil, err
	}

	if _, ok := r.subPlans[s.PlanCode]; !ok {
		return nil, apperrors.NewInternalErr(fmt.Sprintf("subscription has invalid plan code: %v", s.PlanCode))
	}

	return s, nil
}

func (r pgRepo) CreateEvent(se *subPkg.SubscriptionEvent) error {
	sql := `
		INSERT INTO subscription_events (
			active,
			type,
			status,
			subscription_origin,
			subscription_plan_code,
			subscription_period_unit,
			subscription_period_count,
			subscription_period_start_date,
			subscription_period_finish_date,
			user_id,
			subscription_id,
			invoice_uuid,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`
	err := r.pool.QueryRow(
		r.ctx,
		sql,
		se.Active,
		se.Type,
		se.Status,
		se.SubscriptionOrigin,
		se.SubscriptionPlanCode,
		se.SubscriptionPeriodUnit,
		se.SubscriptionPeriodCount,
		se.SubscriptionPeriodStartDate,
		se.SubscriptionPeriodFinishDate,
		se.UserID,
		se.SubscriptionID,
		se.InvoiceUUID,
		se.CreatedAt,
		se.UpdatedAt,
	).Scan(&se.ID)

	return err
}

func (r pgRepo) UpdateEvent(se *subPkg.SubscriptionEvent) error {
	sql := `
		UPDATE subscription_events SET
			active = $1,
			type = $2,
			status = $3,
			subscription_origin = $4,
			subscription_plan_code = $5,
			subscription_period_unit = $6,
			subscription_period_count = $7,
			subscription_period_start_date = $8,
			subscription_period_finish_date = $9,
			user_id = $10,
			subscription_id = $11,
			invoice_uuid = $12,
			created_at = $13,
			updated_at = $14
		WHERE id = $15
	`
	_, err := r.pool.Exec(
		r.ctx,
		sql,
		se.Active,
		se.Type,
		se.Status,
		se.SubscriptionOrigin,
		se.SubscriptionPlanCode,
		se.SubscriptionPeriodUnit,
		se.SubscriptionPeriodCount,
		se.SubscriptionPeriodStartDate,
		se.SubscriptionPeriodFinishDate,
		se.UserID,
		se.SubscriptionID,
		se.InvoiceUUID,
		se.CreatedAt,
		se.UpdatedAt,
		se.ID,
	)

	return err
}

func (r pgRepo) GetActiveEventByInvoiceUUIDAndStatuses(invoiceUUID string, statuses []subPkg.SubscriptionEventStatus) (*subPkg.SubscriptionEvent, error) {
	se := &subPkg.SubscriptionEvent{}

	sql := `
		SELECT
			id,
			active,
			type,
			status,
			subscription_origin,
			subscription_plan_code,
			subscription_period_unit,
			subscription_period_count,
			subscription_period_start_date,
			subscription_period_finish_date,
			user_id,
			subscription_id,
			invoice_uuid,
			created_at,
			updated_at
		FROM subscription_events
		WHERE 
			active = true AND 
			invoice_uuid = $1 AND 
			status = ANY($2)
		LIMIT 1
	`
	err := r.pool.QueryRow(
		r.ctx,
		sql,
		invoiceUUID,
		statuses,
	).Scan(
		&se.ID,
		&se.Active,
		&se.Type,
		&se.Status,
		&se.SubscriptionOrigin,
		&se.SubscriptionPlanCode,
		&se.SubscriptionPeriodUnit,
		&se.SubscriptionPeriodCount,
		&se.SubscriptionPeriodStartDate,
		&se.SubscriptionPeriodFinishDate,
		&se.UserID,
		&se.SubscriptionID,
		&se.InvoiceUUID,
		&se.CreatedAt,
		&se.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.NewNotFoundErr(fmt.Sprintf("couldn't find subscription event by invoice UUID %v and statuses %v", invoiceUUID, statuses))
		}
		return nil, err
	}

	if _, ok := subPkg.SubscriptionEventTypeMapping[string(se.Type)]; !ok {
		return nil, apperrors.NewInternalErr(fmt.Sprintf("subscription event has invalid type: %v", se.Type))
	}
	if _, ok := subPkg.SubscriptionEventStatusMapping[string(se.Status)]; !ok {
		return nil, apperrors.NewInternalErr(fmt.Sprintf("subscription event has invalid status: %v", se.Status))
	}
	if _, ok := subPkg.SubscriptionOriginMapping[string(se.SubscriptionOrigin)]; !ok {
		return nil, apperrors.NewInternalErr(fmt.Sprintf("subscription event has invalid subscription origin: %v", se.SubscriptionOrigin))
	}
	if _, ok := r.subPlans[string(se.SubscriptionPlanCode)]; !ok {
		return nil, apperrors.NewInternalErr(fmt.Sprintf("subscription event has invalid subscription plan code: %v", se.SubscriptionPlanCode))
	}
	if _, ok := subPkg.SubscriptionPeriodUnitMapping[string(se.SubscriptionPeriodUnit)]; !ok {
		return nil, apperrors.NewInternalErr(fmt.Sprintf("subscription event has invalid subscription period unit: %v", se.SubscriptionPeriodUnit))
	}

	return se, nil
}
