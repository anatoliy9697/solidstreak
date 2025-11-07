package repo

import (
	"context"

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
			INSERT INTO habits (active, archived, title, description, color, creator_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, creator_id
		)
		INSERT INTO users_habits (active, user_id, habit_id, is_public)
		SELECT TRUE, habit.creator_id, habit.id, $9 FROM habit
		RETURNING habit_id
	`
	err := r.p.QueryRow(
		r.c,
		sql,
		h.Active,
		h.Archived,
		h.Title,
		h.Description,
		h.Color,
		h.CreatorID,
		h.CreatedAt,
		h.UpdatedAt,
		h.IsPublic,
	).Scan(&h.ID)

	return err
}

func (r pgRepo) Update(h *hPkg.Habit) error {
	sql := `
		WITH habit AS (
			UPDATE habits SET
				active = $1,
				archived = $2,
				title = $3,
				description = $4,
				color = $5,
				updated_at = $6
			WHERE id = $7
			RETURNING id, creator_id
		)
		UPDATE users_habits SET
			is_public = $8
		FROM habit h
		WHERE 
			users_habits.habit_id = h.id
			AND users_habits.user_id = h.creator_id
	`
	_, err := r.p.Exec(
		r.c,
		sql,
		h.Active,
		h.Archived,
		h.Title,
		h.Description,
		h.Color,
		h.UpdatedAt,
		h.ID,
		h.IsPublic,
	)

	return err
}

func (r pgRepo) GetByOwnerIDAndStatus(ownerID int64, status hPkg.HabitStatus, requestedByOwner bool) ([]*hPkg.Habit, error) {
	sql := `
		SELECT h.id, h.active, h.archived, h.title, h.description, h.color, h.creator_id, uh.is_public, h.created_at, h.updated_at
		FROM habits h
		JOIN users_habits uh ON 
			h.id = uh.habit_id 
			AND uh.active = TRUE
		WHERE 
			h.active IS TRUE
			AND uh.user_id = $1
	`
	switch status {
	case hPkg.Active:
		sql += " AND h.archived = FALSE"
	case hPkg.Archived:
		sql += " AND h.archived = TRUE"
	}
	if !requestedByOwner {
		sql += " AND uh.is_public = TRUE"
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
			&h.Archived,
			&h.Title,
			&h.Description,
			&h.Color,
			&h.CreatorID,
			&h.IsPublic,
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

func (r pgRepo) GetByIDAndOwnerID(id int64, ownerID int64, requestedByOwner bool) (*hPkg.Habit, error) {
	sql := `
		WITH habit AS (
			SELECT *
			FROM habits h
			WHERE h.id = $1
		)
		SELECT h.id, h.active, h.archived, h.title, h.description, h.color, h.creator_id, uh.is_public, h.created_at, h.updated_at
		FROM habit h
		JOIN users_habits uh ON 
			h.id = uh.habit_id 
			AND uh.active = TRUE
		WHERE 
			uh.user_id = $2
	`
	if !requestedByOwner {
		sql += " AND uh.is_public = TRUE"
	}

	h := &hPkg.Habit{}
	err := r.p.QueryRow(
		r.c,
		sql,
		id,
		ownerID,
	).Scan(
		&h.ID,
		&h.Active,
		&h.Archived,
		&h.Title,
		&h.Description,
		&h.Color,
		&h.CreatorID,
		&h.IsPublic,
		&h.CreatedAt,
		&h.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound("couldn't find habit for specified user")
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

func (r pgRepo) GetUserHabitsCompletedChecks(userID int64, habitIDs []int64, from, to *date.Date) ([]*hPkg.HabitCheck, error) {
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
