package invoice

import "time"

type Currency string

const (
	CurrencyTgStars Currency = "XTR"
)

var CurrencyMapping = map[string]Currency{
	string(CurrencyTgStars): CurrencyTgStars,
}

type InvoiceStatus string

const (
	InvoiceStatusPending InvoiceStatus = "pending"
	InvoiceStatusPaid    InvoiceStatus = "paid"
	InvoiceStatusExpired InvoiceStatus = "expired"
)

var InvoiceStatusMapping = map[string]InvoiceStatus{
	string(InvoiceStatusPending): InvoiceStatusPending,
	string(InvoiceStatusPaid):    InvoiceStatusPaid,
	string(InvoiceStatusExpired): InvoiceStatusExpired,
}

type Invoice struct {
	Active            bool          `json:"active"`
	UUID              string        `json:"uuid"`
	Status            InvoiceStatus `json:"status"`
	Currency          Currency      `json:"currency"`
	Amount            int64         `json:"amount"`
	UserID            int64         `json:"userId"`
	TgMessageID       int64         `json:"tgMessageId"`
	TgPaymentChargeID string        `json:"tgPaymentChargeId"`
	ExpiresAt         time.Time     `json:"expiresAt"`
	CreatedAt         time.Time     `json:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt"`
}

func NewInvoice(
	uuid string,
	currency Currency,
	amount int64,
	userID int64,
	expiresAt time.Time,
) *Invoice {
	return &Invoice{
		Active:    true,
		UUID:      uuid,
		Status:    InvoiceStatusPending,
		Currency:  currency,
		Amount:    amount,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
