package subscription

import (
	"time"
)

type Pricing struct {
	TgStarsPerMonth float64 `json:"tgStarsPerMonth"`
	TgStarsPerYear  float64 `json:"tgStarsPerYear"`
	TgStarsForever  float64 `json:"tgStarsForever"`
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
