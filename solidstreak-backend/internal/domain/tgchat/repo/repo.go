package repo

import (
	"context"

	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	IsExistsByTgID(int64) (bool, error)
	Create(*tcPkg.Chat) error
	Update(*tcPkg.Chat) error
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
