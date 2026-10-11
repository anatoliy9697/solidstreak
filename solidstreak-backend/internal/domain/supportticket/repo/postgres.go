package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"

	stPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/supportticket"
)

type pgRepo struct {
	c context.Context
	p *pgxpool.Pool
}

func initPGRepo(c context.Context, p *pgxpool.Pool) *pgRepo {
	return &pgRepo{c, p}
}

func (r pgRepo) Create(t *stPkg.Ticket) error {
	sql := `
		INSERT INTO support_tickets (
			active, 
			uuid, 
			status, 
			user_id, 
			tg_chat_id,
			ticket_message, 
			ticket_message_tg_id, 
			ticket_answer, 
			ticket_answer_tg_id, 
			created_at, 
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.p.Exec(r.c, sql,
		t.Active,
		t.UUID,
		t.Status,
		t.UserID,
		t.TgChatID,
		t.TicketMessage,
		t.TicketMessageTgID,
		t.TicketAnswer,
		t.TicketAnswerTgID,
		t.CreatedAt,
		t.UpdatedAt,
	)

	return err
}

func (r pgRepo) Update(t *stPkg.Ticket) error {
	sql := `
		UPDATE support_tickets SET
			active = $1,
			status = $2,
			user_id = $3,
			tg_chat_id = $4,
			ticket_message = $5,
			ticket_message_tg_id = $6,
			ticket_answer = $7,
			ticket_answer_tg_id = $8,
			created_at = $9,
			updated_at = $10
		WHERE uuid = $11
	`
	_, err := r.p.Exec(r.c, sql,
		t.Active,
		t.Status,
		t.UserID,
		t.TgChatID,
		t.TicketMessage,
		t.TicketMessageTgID,
		t.TicketAnswer,
		t.TicketAnswerTgID,
		t.CreatedAt,
		t.UpdatedAt,
		t.UUID,
	)

	return err
}

func (r pgRepo) GetActiveByUserIDAndStatuses(userID int64, statuses []stPkg.TicketStatus) (*stPkg.Ticket, error) {
	t := &stPkg.Ticket{}

	sql := `
		SELECT
			active,
			uuid,
			status,
			user_id,
			tg_chat_id,
			ticket_message,
			ticket_message_tg_id,
			ticket_answer,
			ticket_answer_tg_id,
			created_at,
			updated_at
		FROM support_tickets
		WHERE 
			active = true AND 
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
		&t.Active,
		&t.UUID,
		&t.Status,
		&t.UserID,
		&t.TgChatID,
		&t.TicketMessage,
		&t.TicketMessageTgID,
		&t.TicketAnswer,
		&t.TicketAnswerTgID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.NewNotFoundErr(fmt.Sprintf("couldn't find active support ticket by user ID %v and statuses %v", userID, statuses))
		}
		return nil, err
	}

	if _, ok := stPkg.TicketStatusMapping[string(t.Status)]; !ok {
		return nil, apperrors.NewInternalErr(fmt.Sprintf("ticket has invalid status: %v", t.Status))
	}

	return t, nil
}

func (r pgRepo) GetActiveByUUIDAndStatuses(uuid string, statuses []stPkg.TicketStatus) (*stPkg.Ticket, error) {
	t := &stPkg.Ticket{}

	sql := `
		SELECT
			active,
			uuid,
			status,
			user_id,
			tg_chat_id,
			ticket_message,
			ticket_message_tg_id,
			ticket_answer,
			ticket_answer_tg_id,
			created_at,
			updated_at
		FROM support_tickets
		WHERE 
			active = true AND 
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
		&t.Active,
		&t.UUID,
		&t.Status,
		&t.UserID,
		&t.TgChatID,
		&t.TicketMessage,
		&t.TicketMessageTgID,
		&t.TicketAnswer,
		&t.TicketAnswerTgID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.NewNotFoundErr(fmt.Sprintf("couldn't find active support ticket by UUID %v and statuses %v", uuid, statuses))
		}
		return nil, err
	}

	if _, ok := stPkg.TicketStatusMapping[string(t.Status)]; !ok {
		return nil, apperrors.NewInternalErr(fmt.Sprintf("ticket has invalid status: %v", t.Status))
	}

	return t, nil
}
