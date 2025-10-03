package repo

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	hPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/habit"
	"github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/date"
	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"
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

func (r pgRepo) GetByOwnerID(ownerID int64, gettingMode string) ([]*hPkg.Habit, error) {
	sql := `
		SELECT h.id, h.active, h.title, h.description, h.creator_id, h.created_at, h.updated_at
		FROM habits h
		JOIN users_habits uh ON 
			h.id = uh.habit_id 
			AND uh.active = TRUE
		WHERE 
			uh.user_id = $1
	`
	switch gettingMode {
	case Active:
		sql += " AND h.active = TRUE"
	case NotActive:
		sql += " AND h.active = FALSE"
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

func (r pgRepo) SetUserHabitCheck(hc *hPkg.HabitCheck) error {
	sql := `
		INSERT INTO user_habit_checks (user_id, habit_id, check_date, completed, checked_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (habit_id, user_id, check_date) DO UPDATE SET
			completed = EXCLUDED.completed,
			checked_at = EXCLUDED.checked_at
	`
	_, err := r.p.Exec(
		r.c,
		sql,
		hc.UserID,
		hc.HabitID,
		hc.CheckDate,
		hc.Completed,
		hc.CheckedAt,
	)

	return err
}

func (r pgRepo) GetUserHabitsCompletedChecks(userID int64, habitIDs []int64, from, to date.Date) ([]*hPkg.HabitCheck, error) {
	sql := `
		SELECT habit_id, user_id, completed, check_date, checked_at
		FROM user_habit_checks
		WHERE user_id = $1
			AND habit_id = ANY($2)
			AND check_date >= $3
			AND check_date <= $4
			AND completed = TRUE
		ORDER BY check_date ASC
	`
	rows, err := r.p.Query(
		r.c,
		sql,
		userID,
		habitIDs,
		from,
		to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checks := []*hPkg.HabitCheck{}
	for rows.Next() {
		hc := &hPkg.HabitCheck{}
		err = rows.Scan(
			&hc.HabitID,
			&hc.UserID,
			&hc.Completed,
			&hc.CheckDate,
			&hc.CheckedAt,
		)
		if err != nil {
			return nil, err
		}
		checks = append(checks, hc)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return checks, nil
}
