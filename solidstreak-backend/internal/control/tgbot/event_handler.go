package tgbot

import (
	"context"

	tgbotapi "github.com/mymmrac/telego"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
	usecases "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/usecases"
)

type EventHandler struct {
	Code string
	Res  common.Resources
}

func (eh EventHandler) Run(ctx context.Context, doneCh chan string, upd *tgbotapi.Update) {
	var (
		err error
		usr *usrPkg.User
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
			_ = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.GetErrorTgMessageParams(eh.Res, tc, usr))
		}
		doneCh <- eh.Code
	}()

	eh.Res.Logger = eh.Res.Logger.With("handlerCode", eh.Code)

	if usr, err = usecases.MapUserToInnerAndSave(eh.Res, upd.Message.From); err != nil {
		return
	}
	eh.Res.Logger.Debug("user mapped to inner model and saved to DB", "user", usr)

	if tc, err = usecases.MapTgChatToInnerAndSave(eh.Res, upd.Message.Chat, usr); err != nil {
		return
	}
	eh.Res.Logger.Debug("telegram chat mapped to inner model and saved to DB", "tgChat", tc)

	if err = usecases.SendTgMessage(ctx, eh.Res.TgBotAPI, usecases.GetSuccessTgMessageParams(eh.Res, tc, usr)); err != nil {
		return
	}
}
