package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"

	invPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice"
	st "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/schedulertask"
)

type pgRepo struct {
	c context.Context
	p *pgxpool.Pool
}

func initPGRepo(c context.Context, p *pgxpool.Pool) *pgRepo {
	return &pgRepo{c, p}
}

func (r pgRepo) FetchSchedulerTasksWithLocking(batchSize int, lock_duration time.Duration, lock_owner_id string) ([]st.Task, error) {
	sql := `
		WITH batch AS (
			SELECT
				uuid
			FROM invoices
			WHERE 
				active = true AND 
				status = 'pending' AND 
				expires_at <= NOW() AND 
				(lock_owner_id IS NULL OR locked_at <= NOW() - $3 * INTERVAL '1 millisecond')
			ORDER BY expires_at ASC
			FOR UPDATE SKIP LOCKED
			LIMIT $1
		)
		UPDATE invoices
		SET 
			lock_owner_id = $2,
			locked_at = NOW()
		WHERE uuid IN (SELECT uuid FROM batch)
		RETURNING
			active,
			uuid,
			status,
			currency,
			amount,
			user_id,
			tg_chat_id,
			tg_message_id,
			tg_payment_charge_id,
			expires_at,
			created_at,
			updated_at
	`
	rows, err := r.p.Query(r.c, sql, batchSize, lock_owner_id, int(lock_duration.Milliseconds()))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []st.Task{}
	for rows.Next() {
		i := &invPkg.Invoice{}
		err = rows.Scan(
			&i.Active,
			&i.UUID,
			&i.Status,
			&i.Currency,
			&i.Amount,
			&i.UserID,
			&i.TgChatID,
			&i.TgMessageID,
			&i.TgPaymentChargeID,
			&i.ExpiresAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if _, ok := invPkg.InvoiceStatusMapping[string(i.Status)]; !ok {
			return nil, apperrors.NewInternalErr(fmt.Sprintf("invoice has invalid status: %v", i.Status))
		}
		if _, ok := invPkg.CurrencyMapping[string(i.Currency)]; !ok {
			return nil, apperrors.NewInternalErr(fmt.Sprintf("invoice has invalid currency: %v", i.Currency))
		}
		tasks = append(tasks, invPkg.ProcessExpiredInvoiceTask{Invoice: i})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return tasks, nil
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
			tg_chat_id,
			tg_message_id, 
			tg_payment_charge_id, 
			expires_at, 
			created_at, 
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.p.Exec(r.c, sql,
		i.Active,
		i.UUID,
		i.Status,
		i.Currency,
		i.Amount,
		i.UserID,
		i.TgChatID,
		i.TgMessageID,
		i.TgPaymentChargeID,
		i.ExpiresAt,
		i.CreatedAt,
		i.UpdatedAt,
	)

	return err
}

func (r pgRepo) Update(i *invPkg.Invoice) error {
	sql := `
		UPDATE invoices SET
			active = $1,
			status = $2,
			currency = $3,
			amount = $4,
			user_id = $5,
			tg_chat_id = $6,
			tg_message_id = $7,
			tg_payment_charge_id = $8,
			expires_at = $9,
			lock_owner_id = NULL,
			locked_at = NULL,
			created_at = $10,
			updated_at = $11
		WHERE uuid = $12
	`
	_, err := r.p.Exec(r.c, sql,
		i.Active,
		i.Status,
		i.Currency,
		i.Amount,
		i.UserID,
		i.TgChatID,
		i.TgMessageID,
		i.TgPaymentChargeID,
		i.ExpiresAt,
		i.CreatedAt,
		i.UpdatedAt,
		i.UUID,
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
			tg_chat_id,
			tg_message_id,
			tg_payment_charge_id,
			expires_at,
			created_at,
			updated_at
		FROM invoices
		WHERE 
			active = true AND 
			expires_at > NOW() AND 
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
		&i.TgChatID,
		&i.TgMessageID,
		&i.TgPaymentChargeID,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.NewNotFoundErr(fmt.Sprintf("couldn't find active and not expired invoice by user ID %v and statuses %v", userID, statuses))
		}
		return nil, err
	}

	if _, ok := invPkg.InvoiceStatusMapping[string(i.Status)]; !ok {
		return nil, apperrors.NewInternalErr(fmt.Sprintf("invoice has invalid status: %v", i.Status))
	}

	if _, ok := invPkg.CurrencyMapping[string(i.Currency)]; !ok {
		return nil, apperrors.NewInternalErr(fmt.Sprintf("invoice has invalid currency: %v", i.Currency))
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
			tg_chat_id,
			tg_message_id,
			tg_payment_charge_id,
			expires_at,
			created_at,
			updated_at
		FROM invoices
		WHERE 
			active = true AND 
			expires_at > NOW() AND 
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
		&i.TgChatID,
		&i.TgMessageID,
		&i.TgPaymentChargeID,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.NewNotFoundErr(fmt.Sprintf("couldn't find active and not expired invoice by uuid %v and statuses %v", uuid, statuses))
		}
		return nil, err
	}

	if _, ok := invPkg.InvoiceStatusMapping[string(i.Status)]; !ok {
		return nil, apperrors.NewInternalErr(fmt.Sprintf("invoice has invalid status: %v", i.Status))
	}

	if _, ok := invPkg.CurrencyMapping[string(i.Currency)]; !ok {
		return nil, apperrors.NewInternalErr(fmt.Sprintf("invoice has invalid currency: %v", i.Currency))
	}

	return i, nil
}
