package repo

import (
	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"
)

func (r pgRepo) GetBasic() (*subPkg.Subscription, error) {
	return &subPkg.Subscription{
		Active:   true,
		PlanCode: "basic",
	}, nil
}

func (r pgRepo) GetPlanByCode(planCode string) (*subPkg.Plan, error) {
	plan, ok := r.subPlans[planCode]
	if !ok {
		return nil, apperrors.NewInternalErr("subscription plan not found")
	}

	return plan, nil
}
