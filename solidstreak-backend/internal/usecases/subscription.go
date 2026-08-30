package usecases

import (
	"errors"
	"time"

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

func CalculateSubscriptionEventStartAndFinishDt(renew bool, prevEventFinishDt *time.Time, subscriptionEvent *subPkg.SubscriptionEvent) (*time.Time, *time.Time) {
	var startDt time.Time
	if renew && prevEventFinishDt != nil {
		startDt = time.Date(prevEventFinishDt.Year(), prevEventFinishDt.Month(), prevEventFinishDt.Day(), 0, 0, 0, 0, prevEventFinishDt.Location()).AddDate(0, 0, 1)
	} else {
		startDt = time.Now()
		startDt = time.Date(startDt.Year(), startDt.Month(), startDt.Day(), 0, 0, 0, 0, startDt.Location())
	}

	var finishDtPtr *time.Time
	if subscriptionEvent.SubscriptionPeriodUnit != subPkg.SubscriptionPeriodUnitLifetime {
		var finishDt time.Time
		switch subscriptionEvent.SubscriptionPeriodUnit {
		case subPkg.SubscriptionPeriodUnitYear:
			finishDt = startDt.AddDate(int(subscriptionEvent.SubscriptionPeriodCount), 0, 0)
		case subPkg.SubscriptionPeriodUnitMonth:
			finishDt = startDt.AddDate(0, int(subscriptionEvent.SubscriptionPeriodCount), 0)
		}
		finishDtPtr = &finishDt
	}

	return &startDt, finishDtPtr
}

func CreateOrUpdateSubscription(
	r common.Resources,
	subscription *subPkg.Subscription,
	renew bool,
	planCode string,
	startDt *time.Time,
	finishDt *time.Time,
	userID int64,
) (*subPkg.Subscription, error) {
	var err error

	if renew {
		subscription.SetFinishDt(finishDt)
		err = r.SubRepo.Update(subscription)
	} else {
		subscription = subPkg.NewSubscription(planCode, startDt, finishDt, userID)
		err = r.SubRepo.Create(subscription)
	}

	return subscription, err
}

func CompleteSubscriptionEvent(r common.Resources, subscriptionEvent *subPkg.SubscriptionEvent, startDt *time.Time, finishDt *time.Time, subscriptionID *int64) error {
	subscriptionEvent.SetSubscriptionStartAndFinishDt(startDt, finishDt)
	subscriptionEvent.SetSubscriptionID(subscriptionID)
	subscriptionEvent.MarkAsCompleted()
	return r.SubRepo.UpdateEvent(subscriptionEvent)
}
