package supportticket

import "time"

type TicketStatus string

const (
	TicketStatusNew    TicketStatus = "new"
	TicketStatusClosed TicketStatus = "closed"
)

var TicketStatusAll = []TicketStatus{
	TicketStatusNew,
	TicketStatusClosed,
}

var TicketStatusMapping = map[string]TicketStatus{
	string(TicketStatusNew):    TicketStatusNew,
	string(TicketStatusClosed): TicketStatusClosed,
}

type Ticket struct {
	Active            bool
	UUID              string
	Status            TicketStatus
	UserID            int64
	TgChatID          int64
	TicketMessage     string
	TicketMessageTgID int
	TicketAnswer      *string
	TicketAnswerTgID  *int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewTicket(
	uuid string,
	userID int64,
	tgChatID int64,
	ticketMessage string,
	ticketMessageTgID int,
) *Ticket {
	return &Ticket{
		Active:            true,
		UUID:              uuid,
		Status:            TicketStatusNew,
		UserID:            userID,
		TgChatID:          tgChatID,
		TicketMessage:     ticketMessage,
		TicketMessageTgID: ticketMessageTgID,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}
}
