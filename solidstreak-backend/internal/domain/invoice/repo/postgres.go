package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

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
