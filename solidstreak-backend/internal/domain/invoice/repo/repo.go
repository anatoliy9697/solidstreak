package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	invPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice"
	st "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/schedulertask"
)

type Repo interface {
	FetchSchedulerTasksWithLocking(int, time.Duration, string) ([]st.Task, error)
	Create(*invPkg.Invoice) error
	Update(*invPkg.Invoice) error
	GetActiveNotExpiredByUserIDAndStatuses(int64, []invPkg.InvoiceStatus) (*invPkg.Invoice, error)
	GetActiveNotExpiredByUUIDAndStatuses(string, []invPkg.InvoiceStatus) (*invPkg.Invoice, error)
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
