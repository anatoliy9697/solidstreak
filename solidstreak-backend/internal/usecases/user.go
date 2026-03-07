package usecases

import (
	tgbotapi "github.com/mymmrac/telego"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
)

func mapUserToInner(u *tgbotapi.User) *usrPkg.User {
	return usrPkg.NewUser(
		u.ID,
		u.Username,
		u.FirstName,
		u.LastName,
		u.LanguageCode,
		common.ToLocalLang(u.LanguageCode),
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
