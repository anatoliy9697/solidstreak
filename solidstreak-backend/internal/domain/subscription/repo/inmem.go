package repo

import (
	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"
)

func (r pgRepo) GetBasic() (*subPkg.Subscription, error) {
	const BASIC_PLAN_CODE = "basic"

	plan, ok := r.subPlans[BASIC_PLAN_CODE]
	if !ok {
		return nil, apperrors.ErrInternal("plan for basic subscription is not set")
	}

	return &subPkg.Subscription{
		Active:   true,
		PlanCode: BASIC_PLAN_CODE,
		Plan:     plan,
	}, nil
}

func (r pgRepo) GetPlanByCode(planCode string) (*subPkg.Plan, error) {
	plan, ok := r.subPlans[planCode]
	if !ok {
		return nil, apperrors.ErrNotFound("subscription plan not found")
	}

	return plan, nil
}
