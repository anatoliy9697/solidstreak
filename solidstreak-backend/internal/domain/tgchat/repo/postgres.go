package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
)

type pgRepo struct {
	c context.Context
	p *pgxpool.Pool
}

func initPGRepo(c context.Context, p *pgxpool.Pool) *pgRepo {
	return &pgRepo{c, p}
}

func (r pgRepo) IsExistsByUserID(userID int64) (bool, error) {
	exists := false

	sql := `SELECT EXISTS(SELECT 1 FROM tg_chats WHERE user_id = $1)`
	err := r.p.QueryRow(
		r.c,
		sql,
		userID,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r pgRepo) Create(tc *tcPkg.Chat) error {
	sql := `INSERT INTO tg_chats (tg_id, user_id, created_at) VALUES ($1, $2, $3)`
	_, err := r.p.Exec(
		r.c,
		sql,
		tc.TgID,
		tc.UserID,
		tc.CreatedAt,
	)

	return err
}

func (r pgRepo) UpdateByUserID(c *tcPkg.Chat) error {
	sql := `
		UPDATE tg_chats SET
			tg_id = $1,
			created_at = $2
		WHERE user_id = $3
	`
	_, err := r.p.Exec(
		r.c,
		sql,
		c.TgID,
		c.CreatedAt,
		c.UserID,
	)

	return err
}

func (r pgRepo) ByUserID(userID int64) (*tcPkg.Chat, error) {
	var tc tcPkg.Chat

	sql := `SELECT tg_id, user_id, created_at FROM tg_chats WHERE user_id = $1`
	err := r.p.QueryRow(
		r.c,
		sql,
		userID,
	).Scan(
		&tc.TgID,
		&tc.UserID,
		&tc.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &tc, nil
}
