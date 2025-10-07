package repo

import (
	"context"

	hPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/habit"
	"github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/date"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	Create(*hPkg.Habit) error
	Update(*hPkg.Habit) error
	GetByOwnerIDAndStatus(int64, hPkg.HabitStatus, bool) ([]*hPkg.Habit, error)
	GetByIDAndOwnerID(int64, int64, bool) (*hPkg.Habit, error)
	SetUserHabitCheck(*hPkg.HabitCheck) error
	GetUserHabitsCompletedChecks(int64, []int64, *date.Date, *date.Date) ([]*hPkg.HabitCheck, error)
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
