package tgbot

import (
	"context"
	"errors"

	tgbotapi "github.com/mymmrac/telego"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"

	invPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice"
	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
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

	switch {
	case upd.Message != nil && upd.Message.SuccessfulPayment == nil:
		eh.handleMessage(ctx, upd.Message)
	case upd.PreCheckoutQuery != nil:
		eh.handlePreCheckoutQuery(ctx, upd.PreCheckoutQuery)
	case upd.Message != nil && upd.Message.SuccessfulPayment != nil:
		eh.handleSuccessfulPayment(ctx, upd.Message.SuccessfulPayment)
	}
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
				lang = common.GetDefaultLang()
			}

			if err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewErrorTgMessageParams(eh.Res, tc.TgID, lang)); err != nil {
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

	if err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewDefaultSuccessTgMessageParams(eh.Res, tc.TgID, u.TgFirstName, u.LangCode)); err != nil {
		return
	}
}

func (eh EventHandler) handlePreCheckoutQuery(ctx context.Context, pcq *tgbotapi.PreCheckoutQuery) {
	var (
		err error
		u   *usrPkg.User
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
			lang = common.GetDefaultLang()
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

	var invoice *invPkg.Invoice
	if invoice, err = eh.Res.InvRepo.GetActiveNotExpiredByUUIDAndStatuses(invoiceUUID, []invPkg.InvoiceStatus{invPkg.InvoiceStatusPending}); err != nil {
		return
	}

	if u, err = eh.Res.UsrRepo.GetByID(invoice.UserID); err != nil {
		return
	}

	var subscriptionEvent *subPkg.SubscriptionEvent
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
		err error
		tc  *tcPkg.Chat
		u   *usrPkg.User
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

		if !success && tc != nil {
			var lang string
			if u != nil {
				lang = u.LangCode
			} else {
				lang = common.GetDefaultLang()
			}

			if err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewErrorTgMessageParams(eh.Res, tc.TgID, lang)); err != nil {
				eh.Res.Logger.Error("failed to send error message", "error", err)
			}
		}
	}()

	var invoiceUUID string
	if invoiceUUID, err = usecases.ExtractInvoiceUUIDFromPayload(pmt.InvoicePayload); err != nil {
		return
	}

	var invoice *invPkg.Invoice
	if invoice, err = eh.Res.InvRepo.GetActiveNotExpiredByUUIDAndStatuses(invoiceUUID, []invPkg.InvoiceStatus{invPkg.InvoiceStatusPending}); err != nil {
		return
	}

	if u, err = eh.Res.UsrRepo.GetByID(invoice.UserID); err != nil {
		return
	}
	if tc, err = eh.Res.TCRepo.GetByUserID(u.ID); err != nil {
		return
	}

	var subscriptionEvent *subPkg.SubscriptionEvent
	if subscriptionEvent, err = eh.Res.SubRepo.GetActiveEventByInvoiceUUIDAndStatuses(invoice.UUID, []subPkg.SubscriptionEventStatus{subPkg.SubscriptionEventStatusInProgress}); err != nil {
		return
	}

	if pmt.InvoicePayload != usecases.BuildInvoicePayloadString(invoice, subscriptionEvent) {
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

	renew := subscription.PlanCode != "basic"

	startDt, finishDt := usecases.CalculateSubscriptionEventStartAndFinishDate(renew, subscription.FinishDate, subscriptionEvent)

	if subscription, err = usecases.CreateOrUpdateSubscription(eh.Res, subscription, renew, subscriptionEvent.SubscriptionPlanCode, startDt, finishDt, u.ID); err != nil {
		return
	}

	if err = usecases.MarkInvoiceAsPaid(eh.Res, invoice, pmt.TelegramPaymentChargeID); err != nil {
		return
	}

	if err = usecases.CompleteSubscriptionEvent(eh.Res, subscriptionEvent, startDt, finishDt, subscription.ID); err != nil {
		return
	}

	if err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.NewSubscriptionPurchaseSuccessTgMessageParams(eh.Res, tc.TgID, renew, subscriptionEvent, u.LangCode)); err != nil {
		return
	}
}
