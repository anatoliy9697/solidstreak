package invoice

import (
	"time"

	st "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/schedulertask"
)

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
	TgChatID          int64         `json:"-"`
	TgMessageID       int           `json:"-"`
	TgPaymentChargeID string        `json:"-"`
	ExpiresAt         time.Time     `json:"expiresAt"`
	CreatedAt         time.Time     `json:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt"`
}

type ProcessExpiredInvoiceTask struct {
	Invoice *Invoice
}

func NewInvoice(
	uuid string,
	currency Currency,
	amount int64,
	userID int64,
	tgChatID int64,
	expiresAt time.Time,
) *Invoice {
	return &Invoice{
		Active:    true,
		UUID:      uuid,
		Status:    InvoiceStatusPending,
		Currency:  currency,
		Amount:    amount,
		UserID:    userID,
		TgChatID:  tgChatID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (i *Invoice) SetTgMessageID(tgMessageID int) {
	i.TgMessageID = tgMessageID
	i.UpdatedAt = time.Now().UTC()
}

func (i *Invoice) MarkAsPaid(tgPaymentChargeID string) {
	i.Status = InvoiceStatusPaid
	i.TgPaymentChargeID = tgPaymentChargeID
	i.UpdatedAt = time.Now().UTC()
}

func (i *Invoice) SetStatus(status InvoiceStatus) {
	i.Status = status
	i.UpdatedAt = time.Now().UTC()
}

func (peit ProcessExpiredInvoiceTask) Type() st.TaskType {
	return st.TaskTypeProcessExpiredInvoice
}
