package subscription

import (
	"time"
)

type SubscriptionPeriodUnit string

const (
	SubscriptionPeriodUnitMonth    SubscriptionPeriodUnit = "month"
	SubscriptionPeriodUnitYear     SubscriptionPeriodUnit = "year"
	SubscriptionPeriodUnitLifetime SubscriptionPeriodUnit = "lifetime"
)

var SubscriptionPeriodUnitMapping = map[string]SubscriptionPeriodUnit{
	string(SubscriptionPeriodUnitMonth):    SubscriptionPeriodUnitMonth,
	string(SubscriptionPeriodUnitYear):     SubscriptionPeriodUnitYear,
	string(SubscriptionPeriodUnitLifetime): SubscriptionPeriodUnitLifetime,
}

type Currency string

const (
	CurrencyTgStars Currency = "XTR"
)

var CurrencyMapping = map[string]Currency{
	string(CurrencyTgStars): CurrencyTgStars,
}

type Pricing struct {
	PeriodUnit   SubscriptionPeriodUnit `json:"periodUnit"`
	PeriodCount  int64                  `json:"periodCount"`
	Price        float64                `json:"price"`
	Currency     Currency               `json:"currency"`
	DisplayOrder int64                  `json:"displayOrder"`
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
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
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
				PeriodUnit:   SubscriptionPeriodUnitMonth,
				PeriodCount:  1,
				Price:        basicPlanPriceStarsPerMonth,
				Currency:     CurrencyTgStars,
				DisplayOrder: 1,
			},
			{
				PeriodUnit:   SubscriptionPeriodUnitYear,
				PeriodCount:  1,
				Price:        basicPlanPriceStarsPerYear,
				Currency:     CurrencyTgStars,
				DisplayOrder: 2,
			},
			{
				PeriodUnit:   SubscriptionPeriodUnitLifetime,
				PeriodCount:  1,
				Price:        basicPlanPriceStarsLifetime,
				Currency:     CurrencyTgStars,
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
				PeriodUnit:   SubscriptionPeriodUnitMonth,
				PeriodCount:  1,
				Price:        premiumPlanPriceStarsPerMonth,
				Currency:     CurrencyTgStars,
				DisplayOrder: 1,
			},
			{
				PeriodUnit:   SubscriptionPeriodUnitYear,
				PeriodCount:  1,
				Price:        premiumPlanPriceStarsPerYear,
				Currency:     CurrencyTgStars,
				DisplayOrder: 2,
			},
			{
				PeriodUnit:   SubscriptionPeriodUnitLifetime,
				PeriodCount:  1,
				Price:        premiumPlanPriceStarsLifetime,
				Currency:     CurrencyTgStars,
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
