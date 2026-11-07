package tgbot

import (
	"context"
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
	tgbotapi "github.com/mymmrac/telego"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common/i18n"
	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"

	invPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice"
	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
	stPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/supportticket"
	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/usecases"
)

type EventHandler struct {
	Code string
	Res  common.Resources
}

func (eh EventHandler) Run(ctx context.Context, doneCh chan string, upd *tgbotapi.Update) {
	defer func() { doneCh <- eh.Code }()

	eh.Res.Logger = eh.Res.Logger.With("handlerCode", eh.Code)

	eh.Res.Logger.Debug("event handler started")

	switch {
	case upd.PreCheckoutQuery != nil:
		eh.handlePreCheckoutQuery(ctx, upd.PreCheckoutQuery)

	case upd.Message != nil && upd.Message.SuccessfulPayment != nil:
		eh.handleSuccessfulPayment(ctx, upd.Message.SuccessfulPayment)

	case upd.Message != nil && isCmd(upd.Message, "terms"):
		eh.handleTermsCmd(ctx, upd.Message)

	case upd.Message != nil && isCmd(upd.Message, "support"):
		eh.handleSupportCmd(ctx, upd.Message)

	case upd.Message != nil && isCmd(upd.Message, "paysupport"):
		eh.handlePaySupportCmd(ctx, upd.Message)

	case upd.Message != nil:
		eh.handleMessage(ctx, upd.Message)
	}

	eh.Res.Logger.Debug("event handler finished")
}

func (eh EventHandler) handlePreCheckoutQuery(ctx context.Context, pcq *tgbotapi.PreCheckoutQuery) {
	var (
		err               error
		u                 *usrPkg.User
		invoice           *invPkg.Invoice
		subscriptionEvent *subPkg.SubscriptionEvent
	)

	defer func() {
		success := true

		if r := recover(); r != nil {
			success = false
			eh.Res.Logger.Error("panic recovered in EventHandler.handlePreCheckoutQuery", "panic", r)
		}

		if err != nil {
			success = false
			eh.Res.Logger.Error("event handler error", "error", err)
		}

		var lang string
		if u != nil {
			lang = u.LangCode
		} else {
			lang = i18n.GetDefaultLang()
		}

		tgAnswerPCQParams := usecases.NewTgAnswerPreCheckoutQuery(pcq.ID, success, lang)
		if err = usecases.SendTgAnswerPreCheckoutQuery(ctx, eh.Res.TgBotAPI, tgAnswerPCQParams); err != nil {
			eh.Res.Logger.Error("failed to send answer pre-checkout query", "error", err)
		}
	}()

	var invoiceUUID string
	if invoiceUUID, err = usecases.ExtractInvoiceUUIDFromPayload(pcq.InvoicePayload); err != nil {
		return
	}

	if invoice, err = eh.Res.InvRepo.GetActiveByUUIDAndStatuses(invoiceUUID, []invPkg.InvoiceStatus{invPkg.InvoiceStatusPending}, true); err != nil {
		return
	}

	if u, err = eh.Res.UsrRepo.GetByID(invoice.UserID); err != nil {
		return
	}

	if subscriptionEvent, err = eh.Res.SubRepo.GetActiveEventByInvoiceUUIDAndStatuses(invoice.UUID, []subPkg.SubscriptionEventStatus{subPkg.SubscriptionEventStatusInProgress}); err != nil {
		return
	}

	if pcq.InvoicePayload != usecases.BuildInvoicePayloadString(invoice, subscriptionEvent) {
		err = errors.New("invoice payload mismatch")
		return
	}

	var subscription *subPkg.Subscription
	if subscription, err = usecases.GetUserSubscription(eh.Res, u); err != nil {
		return
	}

	if subscription.PlanCode == "premium" && subscription.FinishDate == nil {
		err = errors.New("cannot purchase or renew a premium subscription while user already has an active lifetime subscription")
		return
	}
}

func (eh EventHandler) handleSuccessfulPayment(ctx context.Context, pmt *tgbotapi.SuccessfulPayment) {
	var (
		err                   error
		tc                    *tcPkg.Chat
		u                     *usrPkg.User
		invoice               *invPkg.Invoice
		subscriptionEvent     *subPkg.SubscriptionEvent
		subscription          *subPkg.Subscription
		renew                 bool
		subscriptionProcessed bool = false
	)

	defer func() {
		success := true

		if r := recover(); r != nil {
			success = false
			eh.Res.Logger.Error("panic recovered in EventHandler.handleSuccessfulPayment", "panic", r)
		}

		if err != nil {
			success = false
			eh.Res.Logger.Error("event handler error", "error", err)
		}

		if !success {
			if u != nil {
				if err = usecases.RefundStarPayment(ctx, eh.Res.TgBotAPI, usecases.NewTgRefundStarPaymentParams(u.TgID, pmt.TelegramPaymentChargeID)); err != nil {
					eh.Res.Logger.Error("failed to refund star payment", "error", err)
				}
			}

			if subscription != nil && subscriptionProcessed {
				if err := usecases.RollbackSubscription(eh.Res, subscription, subscriptionEvent); err != nil {
					eh.Res.Logger.Error("failed to rollback subscription", "error", err)
				}
			}

			if invoice != nil {
				invoice.SetStatus(invPkg.InvoiceStatusCancelled)
				if err = eh.Res.InvRepo.Update(invoice); err != nil {
					eh.Res.Logger.Error("failed to mark invoice as cancelled", "error", err)
				}
			}

			if subscriptionEvent != nil {
				subscriptionEvent.SetStatus(subPkg.SubscriptionEventStatusError)
				if err = eh.Res.SubRepo.UpdateEvent(subscriptionEvent); err != nil {
					eh.Res.Logger.Error("failed to mark subscription event as error", "error", err)
				}
			}

			if tc != nil {
				var lang string
				if u != nil {
					lang = u.LangCode
				} else {
					lang = i18n.GetDefaultLang()
				}

				if _, err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewErrorTgMessageParams(eh.Res, tc.TgID, lang)); err != nil {
					eh.Res.Logger.Error("failed to send error message", "error", err)
				}
			}
		}
	}()

	var invoiceUUID string
	if invoiceUUID, err = usecases.ExtractInvoiceUUIDFromPayload(pmt.InvoicePayload); err != nil {
		return
	}

	if invoice, err = eh.Res.InvRepo.GetActiveByUUIDAndStatuses(invoiceUUID, []invPkg.InvoiceStatus{invPkg.InvoiceStatusPending}, true); err != nil {
		return
	}

	if u, err = eh.Res.UsrRepo.GetByID(invoice.UserID); err != nil {
		return
	}
	if tc, err = eh.Res.TCRepo.GetByUserID(u.ID); err != nil {
		return
	}

	if subscriptionEvent, err = eh.Res.SubRepo.GetActiveEventByInvoiceUUIDAndStatuses(invoice.UUID, []subPkg.SubscriptionEventStatus{subPkg.SubscriptionEventStatusInProgress}); err != nil {
		return
	}

	if pmt.InvoicePayload != usecases.BuildInvoicePayloadString(invoice, subscriptionEvent) {
		err = errors.New("invoice payload mismatch")
		return
	}

	if subscription, err = usecases.GetUserSubscription(eh.Res, u); err != nil {
		return
	}

	if subscription.PlanCode == "premium" && subscription.FinishDate == nil {
		err = errors.New("cannot purchase or renew a premium subscription while user already has an active lifetime subscription")
		return
	}

	renew = subscription.PlanCode != "basic"

	startDt, finishDt := usecases.CalculateSubscriptionEventStartAndFinishDate(renew, subscription.FinishDate, subscriptionEvent)

	subscriptionEvent.SetSubscriptionStartAndFinishDate(startDt, finishDt)

	if subscription, err = usecases.CreateOrUpdateSubscription(eh.Res, subscription, renew, subscriptionEvent); err != nil {
		return
	}

	subscriptionProcessed = true

	subscriptionEvent.SetSubscriptionID(subscription.ID)

	if err = usecases.MarkInvoiceAsPaid(eh.Res, invoice, pmt.TelegramPaymentChargeID); err != nil {
		return
	}

	if err = usecases.CompleteSubscriptionEvent(eh.Res, subscriptionEvent); err != nil {
		return
	}

	if _, err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewSubscriptionPurchaseSuccessTgMessageParams(eh.Res, tc.TgID, renew, subscriptionEvent, u.LangCode)); err != nil {
		return
	}
}

func (eh EventHandler) handleTermsCmd(ctx context.Context, msg *tgbotapi.Message) {
	var (
		err error
		tc  *tcPkg.Chat
		u   *usrPkg.User
	)

	defer func() {
		success := true

		if r := recover(); r != nil {
			success = false
			eh.Res.Logger.Error("panic recovered in EventHandler.handleTermsCmd", "panic", r)
		}

		if err != nil {
			success = false
			eh.Res.Logger.Error("event handler error", "error", err)
		}

		if !success && tc != nil {
			var lang string
			if u != nil {
				lang = u.LangCode
			} else {
				lang = i18n.GetDefaultLang()
			}

			if _, err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewErrorTgMessageParams(eh.Res, tc.TgID, lang)); err != nil {
				eh.Res.Logger.Error("failed to send error message", "error", err)
			}
		}
	}()

	if u, err = usecases.MapUserToInnerAndSave(eh.Res, msg.From); err != nil {
		return
	}
	eh.Res.Logger.Debug("user mapped to inner model and saved to DB", "user", u)

	if tc, err = usecases.MapTgChatToInnerAndSave(eh.Res, msg.Chat, u); err != nil {
		return
	}
	eh.Res.Logger.Debug("telegram chat mapped to inner model and saved to DB", "tgChat", tc)

	if _, err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewTermsTgMessageParams(eh.Res, tc.TgID, u.LangCode)); err != nil {
		return
	}
}

func (eh EventHandler) handleSupportCmd(ctx context.Context, msg *tgbotapi.Message) {
	var (
		err error
		tc  *tcPkg.Chat
		u   *usrPkg.User
	)

	defer func() {
		success := true

		if r := recover(); r != nil {
			success = false
			eh.Res.Logger.Error("panic recovered in EventHandler.handleSupportCmd", "panic", r)
		}

		if err != nil {
			success = false
			eh.Res.Logger.Error("event handler error", "error", err)
		}

		if !success && tc != nil {
			var lang string
			if u != nil {
				lang = u.LangCode
			} else {
				lang = i18n.GetDefaultLang()
			}

			if _, err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewErrorTgMessageParams(eh.Res, tc.TgID, lang)); err != nil {
				eh.Res.Logger.Error("failed to send error message", "error", err)
			}
		}
	}()

	if u, err = usecases.MapUserToInnerAndSave(eh.Res, msg.From); err != nil {
		return
	}
	eh.Res.Logger.Debug("user mapped to inner model and saved to DB", "user", u)

	if tc, err = usecases.MapTgChatToInnerAndSave(eh.Res, msg.Chat, u); err != nil {
		return
	}
	eh.Res.Logger.Debug("telegram chat mapped to inner model and saved to DB", "tgChat", tc)

	if _, err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewSupportTgMessageParams(tc.TgID, u.LangCode)); err != nil {
		return
	}
}

func (eh EventHandler) handlePaySupportCmd(ctx context.Context, msg *tgbotapi.Message) {
	var (
		err error
		tc  *tcPkg.Chat
		u   *usrPkg.User
	)

	defer func() {
		success := true

		if r := recover(); r != nil {
			success = false
			eh.Res.Logger.Error("panic recovered in EventHandler.handlePaySupportCmd", "panic", r)
		}

		if err != nil {
			success = false
			eh.Res.Logger.Error("event handler error", "error", err)
		}

		if !success && tc != nil {
			var lang string
			if u != nil {
				lang = u.LangCode
			} else {
				lang = i18n.GetDefaultLang()
			}

			if _, err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewErrorTgMessageParams(eh.Res, tc.TgID, lang)); err != nil {
				eh.Res.Logger.Error("failed to send error message", "error", err)
			}
		}
	}()

	if u, err = usecases.MapUserToInnerAndSave(eh.Res, msg.From); err != nil {
		return
	}
	eh.Res.Logger.Debug("user mapped to inner model and saved to DB", "user", u)

	if tc, err = usecases.MapTgChatToInnerAndSave(eh.Res, msg.Chat, u); err != nil {
		return
	}
	eh.Res.Logger.Debug("telegram chat mapped to inner model and saved to DB", "tgChat", tc)

	var ticket *stPkg.Ticket
	if ticket, err = eh.Res.STRepo.GetActiveByUserIDAndStatuses(u.ID, []stPkg.TicketStatus{stPkg.TicketStatusNew}); err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return
	}

	err = nil

	if ticket != nil {
		_, err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewAlreadyOpenTicketTgMessageParams(tc.TgID, ticket.UUID, u.LangCode))
		return
	}

	ticketMessage := getStringAfterCmd(msg)

	if ticketMessage == "" {
		_, err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewPaySupportTgMessageParams(tc.TgID, u.LangCode))
		return
	}

	ticket = stPkg.NewTicket(
		uuid.NewString(),
		u.ID,
		tc.TgID,
		ticketMessage,
		msg.MessageID,
	)

	if err = eh.Res.STRepo.Create(ticket); err != nil {
		return
	}

	_, err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewPaySupportTicketCreatedTgMessageParams(tc.TgID, ticket.UUID, u.LangCode))
}

func (eh EventHandler) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	var (
		err error
		tc  *tcPkg.Chat
		u   *usrPkg.User
	)

	defer func() {
		success := true

		if r := recover(); r != nil {
			success = false
			eh.Res.Logger.Error("panic recovered in EventHandler.handleMessage", "panic", r)
		}

		if err != nil {
			success = false
			eh.Res.Logger.Error("event handler error", "error", err)
		}

		if !success && tc != nil {
			var lang string
			if u != nil {
				lang = u.LangCode
			} else {
				lang = i18n.GetDefaultLang()
			}

			if _, err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewErrorTgMessageParams(eh.Res, tc.TgID, lang)); err != nil {
				eh.Res.Logger.Error("failed to send error message", "error", err)
			}
		}
	}()

	if u, err = usecases.MapUserToInnerAndSave(eh.Res, msg.From); err != nil {
		return
	}
	eh.Res.Logger.Debug("user mapped to inner model and saved to DB", "user", u)

	if tc, err = usecases.MapTgChatToInnerAndSave(eh.Res, msg.Chat, u); err != nil {
		return
	}
	eh.Res.Logger.Debug("telegram chat mapped to inner model and saved to DB", "tgChat", tc)

	if _, err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewDefaultSuccessTgMessageParams(eh.Res, tc.TgID, u.TgFirstName, u.LangCode)); err != nil {
		return
	}
}

func isCmd(msg *tgbotapi.Message, name string) bool {
	if msg == nil {
		return false
	}

	fields := strings.Fields(msg.Text)
	if len(fields) == 0 {
		return false
	}

	for _, entity := range msg.Entities {
		if entity.Type != tgbotapi.EntityTypeBotCommand || entity.Offset != 0 {
			continue
		}

		command := strings.SplitN(fields[0], "@", 2)[0]
		return command == "/"+name
	}

	return false
}

func getStringAfterCmd(msg *tgbotapi.Message) string {
	if msg == nil {
		return ""
	}

	for _, entity := range msg.Entities {
		if entity.Type != tgbotapi.EntityTypeBotCommand || entity.Offset != 0 {
			continue
		}

		commandEnd := entity.Length
		if commandEnd < 0 || commandEnd > len(msg.Text) {
			return ""
		}

		if separator, separatorSize := utf8.DecodeRuneInString(msg.Text[commandEnd:]); unicode.IsSpace(separator) {
			commandEnd += separatorSize
		}

		return msg.Text[commandEnd:]
	}

	return ""
}
