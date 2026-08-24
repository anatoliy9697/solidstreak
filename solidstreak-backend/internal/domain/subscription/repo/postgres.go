package repo

import (
	"context"

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

func (r pgRepo) GetActiveByUserID(userID int64) (*subPkg.Subscription, error) {
	s := &subPkg.Subscription{}

	sql := `
		SELECT id, active, plan_code, start_dt, finish_dt, created_at, updated_at
		FROM user_subscriptions 
		WHERE 
			active IS TRUE
			AND user_id = $1
			AND start_dt <= NOW()
			AND (finish_dt IS NULL OR finish_dt >= NOW())
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
		&s.StartDt,
		&s.FinishDt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.NewNotFoundErr("couldn't find user subscription")
		}
		return nil, err
	}

	if _, ok := r.subPlans[s.PlanCode]; !ok {
		return nil, apperrors.NewInternalErr("subscription has invalid plan code")
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
			user_id,
			subscription_id,
			invoice_uuid,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
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
		se.UserID,
		se.SubscriptionID,
		se.InvoiceUUID,
		se.CreatedAt,
		se.UpdatedAt,
	).Scan(&se.ID)

	return err
}
