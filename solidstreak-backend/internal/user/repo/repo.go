package repo

import (
	"context"

	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	IsExists(*usrPkg.User) (bool, error)
	Create(*usrPkg.User) error
	UpdateByTgId(*usrPkg.User) error
	ByID(int64) (*usrPkg.User, error)
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
