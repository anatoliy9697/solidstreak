package tgbot

import (
	"context"

	tgbotapi "github.com/mymmrac/telego"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"

	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/usecases"
)

type EventHandler struct {
	Code string
	Res  common.Resources
}

func (eh EventHandler) Run(ctx context.Context, doneCh chan string, upd *tgbotapi.Update) {
	var (
		err error
		u   *usrPkg.User
		tc  *tcPkg.Chat
	)

	defer func() {
		success := true
		if r := recover(); r != nil {
			success = false
			eh.Res.Logger.Error("panic recovered in EventHandler.Run", "panic", r)
		}
		if err != nil {
			success = false
			eh.Res.Logger.Error("event handler error", "error", err)
		}
		if !success && tc != nil {
			if err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.GetErrorTgMessageParams(eh.Res, tc.TgID, u.LangCode)); err != nil {
				eh.Res.Logger.Error("failed to send error message", "error", err)
			}
		}
		doneCh <- eh.Code
	}()

	eh.Res.Logger = eh.Res.Logger.With("handlerCode", eh.Code)

	switch {
	case upd.Message != nil:
		tc, u, err = eh.handleMessage(ctx, upd.Message)
	case upd.PreCheckoutQuery != nil:
		tc, u, err = eh.handlePreCheckoutQuery(ctx, upd.PreCheckoutQuery)
	}
}

func (eh EventHandler) handleMessage(ctx context.Context, msg *tgbotapi.Message) (*tcPkg.Chat, *usrPkg.User, error) {
	var (
		err error
		tc  *tcPkg.Chat
		u   *usrPkg.User
	)

	if u, err = usecases.MapUserToInnerAndSave(eh.Res, msg.From); err != nil {
		return tc, u, err
	}
	eh.Res.Logger.Debug("user mapped to inner model and saved to DB", "user", u)

	if tc, err = usecases.MapTgChatToInnerAndSave(eh.Res, msg.Chat, u); err != nil {
		return tc, u, err
	}
	eh.Res.Logger.Debug("telegram chat mapped to inner model and saved to DB", "tgChat", tc)

	if err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.GetSuccessTgMessageParams(eh.Res, tc.TgID, u)); err != nil {
		return tc, u, err
	}

	return tc, u, nil
}

func (eh EventHandler) handlePreCheckoutQuery(ctx context.Context, pcq *tgbotapi.PreCheckoutQuery) (*tcPkg.Chat, *usrPkg.User, error) {
	var (
		err error
		ok  = true
		tc  *tcPkg.Chat
		u   *usrPkg.User
	)

	defer func() {
		tgAnswerPCQParams := usecases.GetTgAnswerPreCheckoutQuery(pcq.ID, ok, u.LangCode)
		if err = usecases.SendTgAnswerPreCheckoutQuery(ctx, eh.Res.TgBotAPI, tgAnswerPCQParams); err != nil {
			eh.Res.Logger.Error("failed to send answer pre-checkout query", "error", err)
		}
	}()

	if tc, u, err = usecases.ValidatePreCheckoutQuery(eh.Res, pcq.InvoicePayload); err != nil {
		ok = false
		return tc, u, err
	}

	return tc, u, nil
}
