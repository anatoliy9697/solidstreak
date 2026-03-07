package usecases

import (
	"errors"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"
)

func GetUserSubscription(r common.Resources, u *usrPkg.User) (*subPkg.Subscription, error) {
	sub, err := r.SubRepo.GetActiveByUserID(u.ID)
	if errors.Is(err, apperrors.ErrNotFound) {
		sub, err = r.SubRepo.GetBasic()
	}
	if err != nil {
		return nil, err
	}

	if sub.Plan, err = r.SubRepo.GetPlanByCode(sub.PlanCode); err != nil {
		return nil, err
	}

	return sub, nil
}
