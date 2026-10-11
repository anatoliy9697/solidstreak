package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	stPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/supportticket"
)

type Repo interface {
	Create(*stPkg.Ticket) error
	Update(*stPkg.Ticket) error
	GetActiveByUserIDAndStatuses(int64, []stPkg.TicketStatus) (*stPkg.Ticket, error)
	GetActiveByUUIDAndStatuses(string, []stPkg.TicketStatus) (*stPkg.Ticket, error)
}

func Init(c context.Context, p *pgxpool.Pool) Repo {
	return initPGRepo(c, p)
}
