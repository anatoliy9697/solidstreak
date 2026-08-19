package invoice

import "time"

type InvoiceStatus string

const (
	InvoiceStatusPending   InvoiceStatus = "pending"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

var InvoiceStatusMapping = map[string]InvoiceStatus{
	string(InvoiceStatusPending):   InvoiceStatusPending,
	string(InvoiceStatusPaid):      InvoiceStatusPaid,
	string(InvoiceStatusCancelled): InvoiceStatusCancelled,
}

type Invoice struct {
	Active    bool          `json:"active"`
	UUID      string        `json:"uuid"`
	Status    InvoiceStatus `json:"status"`
	UserID    int64         `json:"userId"`
	ExpiresAt time.Time     `json:"expiresAt"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}
