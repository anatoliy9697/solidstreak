package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	invPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice"
)

type Repo interface {
	Create(*invPkg.Invoice) error
	GetActiveNotExpiredByStatusesAndUserID(userID int64, statuses []string) (*invPkg.Invoice, error)
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
