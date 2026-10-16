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

type TicketAction string

const (
	TicketActionAnswer TicketAction = "answer"
)

var TicketActionMapping = map[string]TicketAction{
	string(TicketActionAnswer): TicketActionAnswer,
}

type Ticket struct {
	Active            bool         `json:"active"`
	UUID              string       `json:"uuid"`
	Status            TicketStatus `json:"status"`
	UserID            int64        `json:"userId"`
	TgChatID          int64        `json:"tgChatId"`
	TicketMessage     string       `json:"ticketMessage"`
	TicketMessageTgID int          `json:"ticketMessageTgId"`
	TicketAnswer      *string      `json:"ticketAnswer"`
	TicketAnswerTgID  *int         `json:"ticketAnswerTgId"`
	CreatedAt         time.Time    `json:"createdAt"`
	UpdatedAt         time.Time    `json:"updatedAt"`
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

func (t *Ticket) SetAnswer(message string, messageTgID int) {
	t.TicketAnswer = &message
	t.TicketAnswerTgID = &messageTgID
	t.UpdatedAt = time.Now().UTC()
}

func (t *Ticket) SetStatus(status TicketStatus) {
	t.Status = status
	t.UpdatedAt = time.Now().UTC()
}
