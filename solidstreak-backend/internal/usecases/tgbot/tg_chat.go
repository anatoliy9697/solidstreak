package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common/resources"
	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
)

func MapTgChatToInnerAndSave(r resources.Resources, outerTC *tgbotapi.Chat, u *usrPkg.User) (*tcPkg.Chat, error) {
	var err error

	tc := tcPkg.NewChat(outerTC.ID, u.ID)

	chatExists := false
	if chatExists, err = r.TCRepo.IsExistsByTgID(tc.TgID); err != nil {
		return nil, err
	}

	if chatExists {
		err = r.TCRepo.Update(tc)
	} else {
		err = r.TCRepo.Create(tc)
	}

	return tc, err
}
