package http

import (
	"encoding/json"
	"net/http"

	tgbotapi "github.com/mymmrac/telego"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"

	stPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/supportticket"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/usecases"
)

type SupportTicketAction struct {
	Action  *string `json:"action"`
	Message *string `json:"message"`
}

type PostSupportTicketActionRequest struct {
	Data *SupportTicketAction `json:"data"`
}

type PostSupportTicketActionResponse struct {
	Data *stPkg.Ticket `json:"data"`
}

func (s Server) postSupportTicketAction(w http.ResponseWriter, r *http.Request) {
	var err error

	logger := s.Res.Logger

	// Adding request ID to request context
	reqID, _ := r.Context().Value(ctxKeyRequestID{}).(string)
	if reqID != "" {
		logger = logger.With("requestId", reqID)
	}

	defer func() {
		if err != nil {
			processError(w, logger, err)
		}
	}()

	var ticketUUID string
	ticketUUID, err = getStringFromURLParams(r, "ticketUUID", true)
	if err != nil {
		return
	}

	var req PostSupportTicketActionRequest

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&req); err != nil {
		err = apperrors.NewBadRequestErr("invalid request payload")
		return
	}

	if req.Data == nil {
		err = apperrors.NewBadRequestErr("support ticket action data is required")
		return
	}
	if req.Data.Action == nil || *req.Data.Action == "" {
		err = apperrors.NewBadRequestErr("support ticket action is required")
		return
	}

	var (
		action stPkg.TicketAction
		ok     bool
	)
	if action, ok = stPkg.TicketActionMapping[*req.Data.Action]; !ok {
		err = apperrors.NewBadRequestErr("invalid support ticket action")
		return
	}

	if action == stPkg.TicketActionAnswer && (req.Data.Message == nil || *req.Data.Message == "") {
		err = apperrors.NewBadRequestErr("message is required for answering a support ticket")
		return
	}

	var ticket *stPkg.Ticket
	if ticket, err = s.Res.STRepo.GetActiveByUUIDAndStatuses(ticketUUID, []stPkg.TicketStatus{stPkg.TicketStatusNew}); err != nil {
		return
	}

	var user *usrPkg.User
	if user, err = s.Res.UsrRepo.GetByID(ticket.UserID); err != nil {
		return
	}

	switch action {
	case stPkg.TicketActionAnswer:
		var tgMsg *tgbotapi.Message
		if tgMsg, err = usecases.SendTgMessage(r.Context(), s.Res.TgBotAPI, usecases.NewSupportTicketAnswerTgMessageParams(ticket.TgChatID, ticket.UUID, *req.Data.Message, string(stPkg.TicketStatusClosed), user.LangCode)); err != nil {
			return
		}
		ticket.SetAnswer(*req.Data.Message, tgMsg.MessageID)
	}

	ticket.SetStatus(stPkg.TicketStatusClosed)

	if err = s.Res.STRepo.Update(ticket); err != nil {
		return
	}

	response := PostSupportTicketActionResponse{Data: ticket}

	json.NewEncoder(w).Encode(response)
}
