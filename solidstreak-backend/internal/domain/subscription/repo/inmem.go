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
		return nil, apperrors.NewNotFoundErr("subscription plan not found")
	}

	return plan, nil
}

func (r pgRepo) GetPlans() ([]*subPkg.Plan, error) {
	plans := make([]*subPkg.Plan, 0, len(r.subPlans))

	for _, plan := range r.subPlans {
		plans = append(plans, plan)
	}

	return plans, nil
}
