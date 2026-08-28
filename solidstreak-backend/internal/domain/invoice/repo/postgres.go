package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"

	invPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice"
)

type pgRepo struct {
	c context.Context
	p *pgxpool.Pool
}

func initPGRepo(c context.Context, p *pgxpool.Pool) *pgRepo {
	return &pgRepo{c, p}
}

func (r pgRepo) Create(i *invPkg.Invoice) error {
	sql := `
		INSERT INTO invoices (
			active, 
			uuid, 
			status, 
			currency, 
			amount, 
			user_id, 
			tg_message_id, 
			tg_payment_charge_id, 
			expires_at, 
			created_at, 
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.p.Exec(r.c, sql,
		i.Active,
		i.UUID,
		i.Status,
		i.Currency,
		i.Amount,
		i.UserID,
		i.TgMessageID,
		i.TgPaymentChargeID,
		i.ExpiresAt,
		i.CreatedAt,
		i.UpdatedAt,
	)

	return err
}

func (r pgRepo) GetActiveNotExpiredByUserIDAndStatuses(userID int64, statuses []invPkg.InvoiceStatus) (*invPkg.Invoice, error) {
	i := &invPkg.Invoice{}

	sql := `
		SELECT
			active,
			uuid,
			status,
			currency,
			amount,
			user_id,
			tg_message_id,
			tg_payment_charge_id,
			expires_at,
			created_at,
			updated_at
		FROM invoices
		WHERE 
			active = true AND 
			expires_at > now() AND 
			user_id = $1 AND 
			status = ANY($2) 
		LIMIT 1
	`
	err := r.p.QueryRow(
		r.c,
		sql,
		userID,
		statuses,
	).Scan(
		&i.Active,
		&i.UUID,
		&i.Status,
		&i.Currency,
		&i.Amount,
		&i.UserID,
		&i.TgMessageID,
		&i.TgPaymentChargeID,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.NewNotFoundErr("couldn't find active and not expired invoice by user id and statuses")
		}
		return nil, err
	}

	if _, ok := invPkg.InvoiceStatusMapping[string(i.Status)]; !ok {
		return nil, apperrors.NewInternalErr("invoice has invalid status")
	}

	if _, ok := invPkg.CurrencyMapping[string(i.Currency)]; !ok {
		return nil, apperrors.NewInternalErr("invoice has invalid currency")
	}

	return i, nil
}

func (r pgRepo) GetActiveNotExpiredByUUIDAndStatuses(uuid string, statuses []invPkg.InvoiceStatus) (*invPkg.Invoice, error) {
	i := &invPkg.Invoice{}

	sql := `
		SELECT
			active,
			uuid,
			status,
			currency,
			amount,
			user_id,
			tg_message_id,
			tg_payment_charge_id,
			expires_at,
			created_at,
			updated_at
		FROM invoices
		WHERE 
			active = true AND 
			expires_at > now() AND 
			uuid = $1 AND 
			status = ANY($2) 
		LIMIT 1
	`
	err := r.p.QueryRow(
		r.c,
		sql,
		uuid,
		statuses,
	).Scan(
		&i.Active,
		&i.UUID,
		&i.Status,
		&i.Currency,
		&i.Amount,
		&i.UserID,
		&i.TgMessageID,
		&i.TgPaymentChargeID,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.NewNotFoundErr("couldn't find active and not expired invoice by uuid and statuses")
		}
		return nil, err
	}

	if _, ok := invPkg.InvoiceStatusMapping[string(i.Status)]; !ok {
		return nil, apperrors.NewInternalErr("invoice has invalid status")
	}

	if _, ok := invPkg.CurrencyMapping[string(i.Currency)]; !ok {
		return nil, apperrors.NewInternalErr("invoice has invalid currency")
	}

	return i, nil
}
