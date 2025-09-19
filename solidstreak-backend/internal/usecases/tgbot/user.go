package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
)

func mapUserToInner(u *tgbotapi.User) *usrPkg.User {
	return usrPkg.NewUser(
		u.ID,
		u.UserName,
		u.FirstName,
		u.LastName,
		u.LanguageCode,
		u.IsBot,
	)
}

func MapUserToInnerAndSave(r common.Resources, outerUsr *tgbotapi.User) (u *usrPkg.User, err error) {
	u = mapUserToInner(outerUsr)

	userExists := false
	if userExists, err = r.UsrRepo.IsExists(u); err == nil {
		if userExists {
			err = r.UsrRepo.Update(u)
		} else {
			err = r.UsrRepo.Create(u)
		}
	}

	return u, err
}
