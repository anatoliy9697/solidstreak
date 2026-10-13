package usecases

import (
	"errors"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
	"github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/date"
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

func CalculateSubscriptionEventStartAndFinishDate(renew bool, prevEventFinishDate *date.Date, se *subPkg.SubscriptionEvent) (*date.Date, *date.Date) {
	var startDate date.Date
	if renew && prevEventFinishDate != nil {
		startDate = prevEventFinishDate.AddDate(0, 0, 1)
	} else {
		startDate = date.Today()
	}

	var finishDatePtr *date.Date
	if se.SubscriptionPeriodUnit != subPkg.SubscriptionPeriodUnitLifetime {
		var finishDate date.Date
		switch se.SubscriptionPeriodUnit {
		case subPkg.SubscriptionPeriodUnitYear:
			finishDate = startDate.AddDate(int(se.SubscriptionPeriodCount), 0, 0)
		case subPkg.SubscriptionPeriodUnitMonth:
			finishDate = startDate.AddDate(0, int(se.SubscriptionPeriodCount), 0)
		}
		finishDatePtr = &finishDate
	}

	return &startDate, finishDatePtr
}

func CreateOrUpdateSubscription(
	r common.Resources,
	s *subPkg.Subscription,
	renew bool,
	se *subPkg.SubscriptionEvent,
) (*subPkg.Subscription, error) {
	var err error

	if renew {
		s.SetFinishDt(se.SubscriptionPeriodFinishDate)
		err = r.SubRepo.Update(s)
	} else {
		s = subPkg.NewSubscription(se.SubscriptionPlanCode, se.SubscriptionPeriodStartDate, se.SubscriptionPeriodFinishDate, se.UserID)
		err = r.SubRepo.Create(s)
	}

	return s, err
}

func CompleteSubscriptionEvent(r common.Resources, se *subPkg.SubscriptionEvent) error {
	se.MarkAsCompleted()
	return r.SubRepo.UpdateEvent(se)
}

func RollbackSubscription(r common.Resources, s *subPkg.Subscription, se *subPkg.SubscriptionEvent) error {
	if s.StartDate.String() == se.SubscriptionPeriodStartDate.String() {
		s.Deactivate()
	} else {
		finishDate := se.SubscriptionPeriodStartDate.AddDate(0, 0, -1)
		s.SetFinishDt(&finishDate)
	}

	return r.SubRepo.Update(s)
}
