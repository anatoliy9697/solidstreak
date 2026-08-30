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

func CalculateSubscriptionEventStartAndFinishDate(renew bool, prevEventFinishDate *date.Date, subscriptionEvent *subPkg.SubscriptionEvent) (*date.Date, *date.Date) {
	var startDate date.Date
	if renew && prevEventFinishDate != nil {
		startDate = prevEventFinishDate.AddDate(0, 0, 1)
	} else {
		startDate = date.Today()
	}

	var finishDatePtr *date.Date
	if subscriptionEvent.SubscriptionPeriodUnit != subPkg.SubscriptionPeriodUnitLifetime {
		var finishDate date.Date
		switch subscriptionEvent.SubscriptionPeriodUnit {
		case subPkg.SubscriptionPeriodUnitYear:
			finishDate = startDate.AddDate(int(subscriptionEvent.SubscriptionPeriodCount), 0, 0)
		case subPkg.SubscriptionPeriodUnitMonth:
			finishDate = startDate.AddDate(0, int(subscriptionEvent.SubscriptionPeriodCount), 0)
		}
		finishDatePtr = &finishDate
	}

	return &startDate, finishDatePtr
}

func CreateOrUpdateSubscription(
	r common.Resources,
	subscription *subPkg.Subscription,
	renew bool,
	planCode string,
	startDate *date.Date,
	finishDate *date.Date,
	userID int64,
) (*subPkg.Subscription, error) {
	var err error

	if renew {
		subscription.SetFinishDt(finishDate)
		err = r.SubRepo.Update(subscription)
	} else {
		subscription = subPkg.NewSubscription(planCode, startDate, finishDate, userID)
		err = r.SubRepo.Create(subscription)
	}

	return subscription, err
}

func CompleteSubscriptionEvent(r common.Resources, subscriptionEvent *subPkg.SubscriptionEvent, startDate *date.Date, finishDate *date.Date, subscriptionID *int64) error {
	subscriptionEvent.SetSubscriptionStartAndFinishDate(startDate, finishDate)
	subscriptionEvent.SetSubscriptionID(subscriptionID)
	subscriptionEvent.MarkAsCompleted()
	return r.SubRepo.UpdateEvent(subscriptionEvent)
}
