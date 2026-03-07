package subscription

import (
	"time"
)

type Pricing struct {
	TgStarsPerMonth int `json:"tgStarsPerMonth"`
	TgStarsPerYear  int `json:"tgStarsPerYear"`
	TgStarsForever  int `json:"tgStarsForever"`
}

type Plan struct {
	Code              string  `json:"code"`
	Price             Pricing `json:"price"`
	ActiveHabitsLimit int     `json:"activeHabitsLimit"`
	ShowAds           bool    `json:"showAds"`
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
	basicPlanActiveHabitsLimit int,
	premiumPlanActiveHabitsLimit int,
	premiumPlanPriceStarsPerMonth int,
	premiumPlanPriceStarsPerYear int,
	premiumPlanPriceStarsForever int,
) map[string]*Plan {
	basicPlan := &Plan{
		Code: "basic",
		Price: Pricing{
			TgStarsPerMonth: 0,
			TgStarsPerYear:  0,
			TgStarsForever:  0,
		},
		ActiveHabitsLimit: basicPlanActiveHabitsLimit,
		ShowAds:           true,
	}

	premiumPlan := &Plan{
		Code: "premium",
		Price: Pricing{
			TgStarsPerMonth: premiumPlanPriceStarsPerMonth,
			TgStarsPerYear:  premiumPlanPriceStarsPerYear,
			TgStarsForever:  premiumPlanPriceStarsForever,
		},
		ActiveHabitsLimit: premiumPlanActiveHabitsLimit,
		ShowAds:           false,
	}

	return map[string]*Plan{
		"basic":   basicPlan,
		"premium": premiumPlan,
	}
}
