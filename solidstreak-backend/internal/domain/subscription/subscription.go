package subscription

import (
	"time"
)

type SubscriptionPeriod string

const (
	Month    SubscriptionPeriod = "month"
	Year     SubscriptionPeriod = "year"
	Lifetime SubscriptionPeriod = "lifetime"
)

var SubscriptionPeriodMapping = map[string]SubscriptionPeriod{
	string(Month):    Month,
	string(Year):     Year,
	string(Lifetime): Lifetime,
}

type Currency string

const (
	TgStars Currency = "XTR"
)

var CurrencyMapping = map[string]Currency{
	string(TgStars): TgStars,
}

type Pricing struct {
	Period       SubscriptionPeriod `json:"period"`
	Price        float64            `json:"price"`
	Currency     Currency           `json:"currency"`
	DisplayOrder int                `json:"displayOrder"`
}

type Plan struct {
	Code              string    `json:"code"`
	Pricing           []Pricing `json:"pricing"`
	ActiveHabitsLimit int64     `json:"activeHabitsLimit"`
	ShowAds           bool      `json:"showAds"`
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
	basicPlanPriceStarsPerMonth float64,
	basicPlanPriceStarsPerYear float64,
	basicPlanPriceStarsLifetime float64,
	basicPlanActiveHabitsLimit int64,
	basicPlanShowAds bool,
	premiumPlanPriceStarsPerMonth float64,
	premiumPlanPriceStarsPerYear float64,
	premiumPlanPriceStarsLifetime float64,
	premiumPlanActiveHabitsLimit int64,
	premiumPlanShowAds bool,
) map[string]*Plan {
	basicPlan := &Plan{
		Code: "basic",
		Pricing: []Pricing{
			{
				Period:       Month,
				Price:        basicPlanPriceStarsPerMonth,
				Currency:     TgStars,
				DisplayOrder: 1,
			},
			{
				Period:       Year,
				Price:        basicPlanPriceStarsPerYear,
				Currency:     TgStars,
				DisplayOrder: 2,
			},
			{
				Period:       Lifetime,
				Price:        basicPlanPriceStarsLifetime,
				Currency:     TgStars,
				DisplayOrder: 3,
			},
		},
		ActiveHabitsLimit: basicPlanActiveHabitsLimit,
		ShowAds:           basicPlanShowAds,
	}

	premiumPlan := &Plan{
		Code: "premium",
		Pricing: []Pricing{
			{
				Period:       Month,
				Price:        premiumPlanPriceStarsPerMonth,
				Currency:     TgStars,
				DisplayOrder: 1,
			},
			{
				Period:       Year,
				Price:        premiumPlanPriceStarsPerYear,
				Currency:     TgStars,
				DisplayOrder: 2,
			},
			{
				Period:       Lifetime,
				Price:        premiumPlanPriceStarsLifetime,
				Currency:     TgStars,
				DisplayOrder: 3,
			},
		},
		ActiveHabitsLimit: premiumPlanActiveHabitsLimit,
		ShowAds:           premiumPlanShowAds,
	}

	return map[string]*Plan{
		"basic":   basicPlan,
		"premium": premiumPlan,
	}
}
