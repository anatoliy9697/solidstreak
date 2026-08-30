package subscription

import (
	"time"

	invPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice"
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

type Pricing struct {
	PeriodUnit   SubscriptionPeriodUnit `json:"periodUnit"`
	PeriodCount  int64                  `json:"periodCount"`
	Price        int64                  `json:"price"`
	Currency     invPkg.Currency        `json:"currency"`
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
	UserID    *int64     `json:"userId,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type SubscriptionEventType string

const (
	SubscriptionEventTypeAcquisition SubscriptionEventType = "acquisition"
)

var SubscriptionEventTypeMapping = map[string]SubscriptionEventType{
	string(SubscriptionEventTypeAcquisition): SubscriptionEventTypeAcquisition,
}

type SubscriptionEventStatus string

const (
	SubscriptionEventStatusInProgress     SubscriptionEventStatus = "in_progress"
	SubscriptionEventStatusCompleted      SubscriptionEventStatus = "completed"
	SubscriptionEventStatusPaymentTimeOut SubscriptionEventStatus = "payment_timeout"
)

var SubscriptionEventStatusMapping = map[string]SubscriptionEventStatus{
	string(SubscriptionEventStatusInProgress):     SubscriptionEventStatusInProgress,
	string(SubscriptionEventStatusCompleted):      SubscriptionEventStatusCompleted,
	string(SubscriptionEventStatusPaymentTimeOut): SubscriptionEventStatusPaymentTimeOut,
}

type SubscriptionOrigin string

const (
	SubscriptionOriginPurchase SubscriptionOrigin = "purchase"
)

var SubscriptionOriginMapping = map[string]SubscriptionOrigin{
	string(SubscriptionOriginPurchase): SubscriptionOriginPurchase,
}

type SubscriptionEvent struct {
	ID                         int64
	Active                     bool
	Type                       SubscriptionEventType
	Status                     SubscriptionEventStatus
	SubscriptionOrigin         SubscriptionOrigin
	SubscriptionPlanCode       string
	SubscriptionPeriodUnit     SubscriptionPeriodUnit
	SubscriptionPeriodCount    int64
	SubscriptionPeriodStartDt  *time.Time
	SubscriptionPeriodFinishDt *time.Time
	UserID                     int64
	SubscriptionID             *int64
	InvoiceUUID                string
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

func GetSubscriptionPlans(
	basicPlanPriceStarsPerMonth int64,
	basicPlanPriceStarsPerYear int64,
	basicPlanPriceStarsLifetime int64,
	basicPlanActiveHabitsLimit int64,
	basicPlanShowAds bool,
	premiumPlanPriceStarsPerMonth int64,
	premiumPlanPriceStarsPerYear int64,
	premiumPlanPriceStarsLifetime int64,
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
				Currency:     invPkg.CurrencyTgStars,
				DisplayOrder: 1,
			},
			{
				PeriodUnit:   SubscriptionPeriodUnitYear,
				PeriodCount:  1,
				Price:        basicPlanPriceStarsPerYear,
				Currency:     invPkg.CurrencyTgStars,
				DisplayOrder: 2,
			},
			{
				PeriodUnit:   SubscriptionPeriodUnitLifetime,
				PeriodCount:  1,
				Price:        basicPlanPriceStarsLifetime,
				Currency:     invPkg.CurrencyTgStars,
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
				Currency:     invPkg.CurrencyTgStars,
				DisplayOrder: 1,
			},
			{
				PeriodUnit:   SubscriptionPeriodUnitYear,
				PeriodCount:  1,
				Price:        premiumPlanPriceStarsPerYear,
				Currency:     invPkg.CurrencyTgStars,
				DisplayOrder: 2,
			},
			{
				PeriodUnit:   SubscriptionPeriodUnitLifetime,
				PeriodCount:  1,
				Price:        premiumPlanPriceStarsLifetime,
				Currency:     invPkg.CurrencyTgStars,
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

func NewSubscription(
	planCode string,
	startDt *time.Time,
	finishDt *time.Time,
	userID int64,
) *Subscription {
	now := time.Now()
	return &Subscription{
		Active:    true,
		PlanCode:  planCode,
		StartDt:   startDt,
		FinishDt:  finishDt,
		CreatedAt: &now,
		UpdatedAt: &now,
		UserID:    &userID,
	}
}

func (s *Subscription) SetFinishDt(finishDt *time.Time) {
	now := time.Now()
	s.FinishDt = finishDt
	s.UpdatedAt = &now
}

func NewSubscriptionEvent(
	eventType SubscriptionEventType,
	status SubscriptionEventStatus,
	subscriptionOrigin SubscriptionOrigin,
	subscriptionPlanCode string,
	subscriptionPeriodUnit SubscriptionPeriodUnit,
	subscriptionPeriodCount int64,
	userID int64,
	invoiceUUID string,
) *SubscriptionEvent {
	return &SubscriptionEvent{
		Active:                  true,
		Type:                    eventType,
		Status:                  status,
		SubscriptionOrigin:      subscriptionOrigin,
		SubscriptionPlanCode:    subscriptionPlanCode,
		SubscriptionPeriodUnit:  subscriptionPeriodUnit,
		SubscriptionPeriodCount: subscriptionPeriodCount,
		UserID:                  userID,
		InvoiceUUID:             invoiceUUID,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}
}

func (se *SubscriptionEvent) SetSubscriptionStartAndFinishDt(startDt, finishDt *time.Time) {
	se.SubscriptionPeriodStartDt = startDt
	se.SubscriptionPeriodFinishDt = finishDt
	se.UpdatedAt = time.Now()
}

func (se *SubscriptionEvent) SetSubscriptionID(subscriptionID *int64) {
	se.SubscriptionID = subscriptionID
	se.UpdatedAt = time.Now()
}

func (se *SubscriptionEvent) MarkAsCompleted() {
	se.Status = SubscriptionEventStatusCompleted
	se.UpdatedAt = time.Now()
}
