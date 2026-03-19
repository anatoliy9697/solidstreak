package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
)

type Repo interface {
	GetActiveByUserID(int64) (*subPkg.Subscription, error)
	GetBasic() (*subPkg.Subscription, error)
	GetPlanByCode(string) (*subPkg.Plan, error)
	GetPlans() ([]*subPkg.Plan, error)
}

func Init(c context.Context, p *pgxpool.Pool, subPlans map[string]*subPkg.Plan) Repo {
	return initPGRepo(c, p, subPlans)
}
