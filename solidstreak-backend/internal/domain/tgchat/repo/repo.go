package repo

import (
	"context"

	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	IsExistsByUserID(int64) (bool, error)
	Create(*tcPkg.Chat) error
	UpdateByUserID(*tcPkg.Chat) error
	ByUserID(int64) (*tcPkg.Chat, error)
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
