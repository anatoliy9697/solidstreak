package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
	usecases "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/usecases/tgbot"
)

type EventHandler struct {
	Code string
	Res  common.Resources
}

func (eh EventHandler) Run(doneCh chan string, upd *tgbotapi.Update) {
	var (
		err error
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
			_ = usecases.SendReplyMsg(eh.Res, tc, "Something went wrong\nPlease try again later")
		}
		doneCh <- eh.Code
	}()

	eh.Res.Logger = eh.Res.Logger.With("handlerCode", eh.Code)

	var usr *usrPkg.User
	if usr, err = usecases.MapUserToInnerAndSave(eh.Res, upd.SentFrom()); err != nil {
		return
	}
	eh.Res.Logger.Debug("user mapped to inner model and saved to DB", "user", usr)

	if tc, err = usecases.MapTgChatToInnerAndSave(eh.Res, upd.FromChat(), usr); err != nil {
		return
	}
	eh.Res.Logger.Debug("telegram chat mapped to inner model and saved to DB", "tgChat", tc)

	if err = usecases.SendReplyMsg(eh.Res, tc, "Hello, "+usr.TgFirstName+"!\nPush \"Open\" button to start using bot"); err != nil {
		return
	}
}
