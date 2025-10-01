package repo

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common/errors"
	hPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/habit"
)

type pgRepo struct {
	c context.Context
	p *pgxpool.Pool
}

func initPGRepo(c context.Context, p *pgxpool.Pool) *pgRepo {
	return &pgRepo{c, p}
}

func (r pgRepo) Create(h *hPkg.Habit) error {
	sql := `
		WITH habit AS (
			INSERT INTO habits (active, title, description, creator_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, creator_id
		)
		INSERT INTO users_habits (active, user_id, habit_id)
		SELECT TRUE, habit.creator_id, habit.id FROM habit
		RETURNING habit_id
	`
	err := r.p.QueryRow(
		r.c,
		sql,
		h.Active,
		h.Title,
		h.Description,
		h.CreatorID,
		h.CreatedAt,
		h.UpdatedAt,
	).Scan(&h.ID)

	return err
}

func (r pgRepo) Update(h *hPkg.Habit) error {
	sql := `
		UPDATE habits SET
			active = $1,
			title = $2,
			description = $3,
			updated_at = $4
		WHERE id = $5
	`
	_, err := r.p.Exec(
		r.c,
		sql,
		h.Active,
		h.Title,
		h.Description,
		h.UpdatedAt,
		h.ID,
	)

	return err
}

func (r pgRepo) GetByID(id int64) (*hPkg.Habit, error) {
	sql := `
		SELECT id, active, title, description, creator_id, created_at, updated_at
		FROM habits WHERE id = $1
	`
	h := &hPkg.Habit{}
	err := r.p.QueryRow(
		r.c,
		sql,
		id,
	).Scan(
		&h.ID,
		&h.Active,
		&h.Title,
		&h.Description,
		&h.CreatorID,
		&h.CreatedAt,
		&h.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound("couldn't find habit with id " + strconv.FormatInt(id, 10))
		}
		return nil, err
	}

	return h, nil
}

func (r pgRepo) GetByOwnerID(ownerID int64, onlyActive bool) ([]*hPkg.Habit, error) {
	sql := `
		SELECT h.id, h.active, h.title, h.description, h.creator_id, h.created_at, h.updated_at
		FROM habits h
		JOIN users_habits uh ON 
			h.id = uh.habit_id 
			AND uh.active = TRUE
		WHERE 
			uh.user_id = $1
	`
	if onlyActive {
		sql += " AND h.active = TRUE"
	}

	rows, err := r.p.Query(r.c, sql, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits := []*hPkg.Habit{}
	for rows.Next() {
		h := &hPkg.Habit{}
		err = rows.Scan(
			&h.ID,
			&h.Active,
			&h.Title,
			&h.Description,
			&h.CreatorID,
			&h.CreatedAt,
			&h.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		habits = append(habits, h)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return habits, nil
}

func (r pgRepo) GetByIDAndOwnerID(id, ownerID int64) (*hPkg.Habit, error) {
	sql := `
		WITH habit AS (
			SELECT *
			FROM habits h
			WHERE h.id = $1
		)
		SELECT h.id, h.active, h.title, h.description, h.creator_id, h.created_at, h.updated_at
		FROM habit h
		JOIN users_habits uh ON 
			h.id = uh.habit_id 
			AND uh.active = TRUE
		WHERE 
			uh.user_id = $2
	`

	h := &hPkg.Habit{}
	err := r.p.QueryRow(
		r.c,
		sql,
		id,
		ownerID,
	).Scan(
		&h.ID,
		&h.Active,
		&h.Title,
		&h.Description,
		&h.CreatorID,
		&h.CreatedAt,
		&h.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound("couldn't find habit with id " + strconv.FormatInt(id, 10) + " for specified user")
		}
		return nil, err
	}

	return h, nil
}
