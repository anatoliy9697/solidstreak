package subscription

import (
	"time"
)

type Pricing struct {
	TgStarsPerMonth int64 `json:"tgStarsPerMonth"`
	TgStarsPerYear  int64 `json:"tgStarsPerYear"`
	TgStarsForever  int64 `json:"tgStarsForever"`
}

type Plan struct {
	Code        string  `json:"code"`
	Price       Pricing `json:"price"`
	HabitsLimit int64   `json:"habitsLimit"`
	ShowAds     bool    `json:"showAds"`
}

type Subscription struct {
	ID        *int64     `json:"id,omitempty"`
	Active    bool       `json:"active"`
	PlanCode  string     `json:"planCode"`
	Plan      *Plan      `json:"plan"`
	StartDt   *time.Time `json:"startDt,omitempty"`
	FinishDt  *time.Time `json:"finishDt,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

func GetSubscriptionPlans(
	basicPlanHabitsLimit int64,
	premiumPlanHabitsLimit int64,
	premiumPlanPriceStarsPerMonth int64,
	premiumPlanPriceStarsPerYear int64,
	premiumPlanPriceStarsForever int64,
) map[string]*Plan {
	basicPlan := &Plan{
		Code: "basic",
		Price: Pricing{
			TgStarsPerMonth: 0,
			TgStarsPerYear:  0,
			TgStarsForever:  0,
		},
		HabitsLimit: basicPlanHabitsLimit,
		ShowAds:     true,
	}

	premiumPlan := &Plan{
		Code: "premium",
		Price: Pricing{
			TgStarsPerMonth: premiumPlanPriceStarsPerMonth,
			TgStarsPerYear:  premiumPlanPriceStarsPerYear,
			TgStarsForever:  premiumPlanPriceStarsForever,
		},
		HabitsLimit: premiumPlanHabitsLimit,
		ShowAds:     false,
	}

	return map[string]*Plan{
		"basic":   basicPlan,
		"premium": premiumPlan,
	}
}
