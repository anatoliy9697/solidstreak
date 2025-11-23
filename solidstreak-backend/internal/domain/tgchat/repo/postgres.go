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

func (r pgRepo) IsExistsByTgID(tgID int64) (bool, error) {
	exists := false

	sql := `SELECT EXISTS(SELECT 1 FROM tg_chats WHERE tg_id = $1)`
	err := r.p.QueryRow(
		r.c,
		sql,
		tgID,
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

func (r pgRepo) Update(c *tcPkg.Chat) error {
	sql := `
		UPDATE tg_chats SET
			user_id = $1
		WHERE tg_id = $2
	`
	_, err := r.p.Exec(
		r.c,
		sql,
		c.UserID,
		c.TgID,
	)

	return err
}
