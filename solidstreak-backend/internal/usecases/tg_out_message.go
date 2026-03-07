package usecases

import (
	"context"
	"fmt"

	tgbotapi "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
)

func GetErrorReplyMsg(r common.Resources, tc *tcPkg.Chat, u *usrPkg.User) *tgbotapi.SendMessageParams {
	lang := ""
	if u != nil {
		lang = u.LangCode
	}
	if lang == "" {
		lang = common.GetDefaultLang()
	}

	return tu.Message(
		tu.ID(tc.TgID),
		common.MESSAGES[lang]["smthWrong"],
	)
}

func GetSuccessReplyMsg(r common.Resources, tc *tcPkg.Chat, u *usrPkg.User) *tgbotapi.SendMessageParams {
	lang := ""
	if u != nil {
		lang = u.LangCode
	}
	if lang == "" {
		lang = common.GetDefaultLang()
	}

	return tu.Message(
		tu.ID(tc.TgID),
		fmt.Sprintf(common.MESSAGES[lang]["helloMsg"], u.TgFirstName),
	).WithReplyMarkup(tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(common.MESSAGES[lang]["open"]).WithWebApp(tu.WebAppInfo(r.WebAppURL)),
		),
	))
}

func SendReplyMsg(ctx context.Context, tgBotAPI *tgbotapi.Bot, msg *tgbotapi.SendMessageParams) error {
	_, err := tgBotAPI.SendMessage(ctx, msg)

	return err
}
