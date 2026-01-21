package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"

	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
)

type pgRepo struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func initPGRepo(c context.Context, p *pgxpool.Pool) *pgRepo {
	return &pgRepo{c, p}
}

func (r pgRepo) IsExists(u *usrPkg.User) (bool, error) {
	exists := false

	sql := `SELECT EXISTS(SELECT 1 FROM users WHERE tg_id = $1)`
	err := r.pool.QueryRow(r.ctx, sql, u.TgID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r pgRepo) Create(u *usrPkg.User) error {
	sql := `
		INSERT INTO users (tg_id, tg_username, tg_first_name, tg_last_name, tg_lang_code, tg_is_bot, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := r.pool.QueryRow(
		r.ctx,
		sql,
		u.TgID,
		u.TgUsername,
		u.TgFirstName,
		u.TgLastName,
		u.TgLangCode,
		u.TgIsBot,
		u.CreatedAt,
	).Scan(&u.ID)

	return err
}

func (r pgRepo) Update(u *usrPkg.User) error {
	sql := `
		UPDATE users SET
			tg_username = $1,
			tg_first_name = $2,
			tg_last_name = $3,
			tg_lang_code = $4,
			tg_is_bot = $5
		WHERE tg_id = $6
		RETURNING id, created_at
	`
	err := r.pool.QueryRow(
		r.ctx,
		sql,
		u.TgUsername,
		u.TgFirstName,
		u.TgLastName,
		u.TgLangCode,
		u.TgIsBot,
		u.TgID,
	).Scan(
		&u.ID,
		&u.CreatedAt,
	)

	return err
}

func (r pgRepo) GetByID(ID int64) (*usrPkg.User, error) {
	u := &usrPkg.User{}

	sql := `
		SELECT id, tg_id, tg_username, tg_first_name, tg_last_name, tg_lang_code, tg_is_bot, created_at
		FROM users WHERE id = $1
	`
	err := r.pool.QueryRow(
		r.ctx,
		sql,
		ID,
	).Scan(
		&u.ID,
		&u.TgID,
		&u.TgUsername,
		&u.TgFirstName,
		&u.TgLastName,
		&u.TgLangCode,
		&u.TgIsBot,
		&u.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound("couldn't find user")
		}
		return nil, err
	}

	return u, nil
}

func (r pgRepo) GetByTgID(tgID int64) (*usrPkg.User, error) {
	u := &usrPkg.User{}

	sql := `
		SELECT id, tg_id, tg_username, tg_first_name, tg_last_name, tg_lang_code, tg_is_bot, created_at
		FROM users WHERE tg_id = $1
	`
	err := r.pool.QueryRow(
		r.ctx,
		sql,
		tgID,
	).Scan(
		&u.ID,
		&u.TgID,
		&u.TgUsername,
		&u.TgFirstName,
		&u.TgLastName,
		&u.TgLangCode,
		&u.TgIsBot,
		&u.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound("couldn't find user")
		}
		return nil, err
	}

	return u, nil
}
