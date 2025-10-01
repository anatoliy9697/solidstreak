package repo

import (
	"context"

	hPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/habit"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	Create(h *hPkg.Habit) error
	Update(h *hPkg.Habit) error
	GetByID(id int64) (*hPkg.Habit, error)
	GetByOwnerID(ownerID int64, onlyActive bool) ([]*hPkg.Habit, error)
	GetByIDAndOwnerID(id, ownerID int64) (*hPkg.Habit, error)
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
