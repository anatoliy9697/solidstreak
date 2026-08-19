package usecases

import (
	"context"
	"fmt"

	tgbotapi "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
)

func GetErrorTgMessageParams(r common.Resources, tc *tcPkg.Chat, u *usrPkg.User) *tgbotapi.SendMessageParams {
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

func GetSuccessTgMessageParams(r common.Resources, tc *tcPkg.Chat, u *usrPkg.User) *tgbotapi.SendMessageParams {
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

func SendTgMessage(ctx context.Context, tgBotAPI *tgbotapi.Bot, msg *tgbotapi.SendMessageParams) error {
	_, err := tgBotAPI.SendMessage(ctx, msg)

	return err
}

func GetTgInvoiceParams(r common.Resources, tc *tcPkg.Chat, u *usrPkg.User, subscriptionPlan *subPkg.Plan, subscriptionPeriodUnit subPkg.SubscriptionPeriodUnit, subscriptionPeriodCount int64, currency subPkg.Currency) (*tgbotapi.SendInvoiceParams, error) {
	lang := ""
	if u != nil {
		lang = u.LangCode
	}
	if lang == "" {
		lang = common.GetDefaultLang()
	}

	var pricing *subPkg.Pricing
	for _, p := range subscriptionPlan.Pricing {
		if p.PeriodUnit == subscriptionPeriodUnit && p.PeriodCount == subscriptionPeriodCount && p.Currency == currency {
			pricing = &p
			break
		}
	}
	if pricing == nil {
		return nil, apperrors.NewNotFoundErr(fmt.Sprintf("couldn't find pricing for plan \"%s\", period unit \"%s\", period count \"%d\", currency \"%s\"", subscriptionPlan.Code, subscriptionPeriodUnit, subscriptionPeriodCount, currency))
	}

	payload := "user:" + fmt.Sprint(u.ID) + ":subscription:" + subscriptionPlan.Code + ":periodUnit:" + string(subscriptionPeriodUnit) + ":periodCount:" + fmt.Sprint(subscriptionPeriodCount) + ":currency:" + string(currency)

	subscriptionPeriodLabel := getSubscriptionPeriodLabel(lang, subscriptionPeriodUnit, subscriptionPeriodCount)
	return &tgbotapi.SendInvoiceParams{
		ChatID:        tgbotapi.ChatID{ID: tc.TgID},
		Title:         common.MESSAGES[lang]["premiumSubscription"],                             // Если появятся другие планы, то нужно будет формировать динамически
		Description:   common.MESSAGES[lang]["accessToPremium"] + " " + subscriptionPeriodLabel, // Если появятся другие планы, то нужно будет формировать динамически + TODO: добавить время, за которое надо выполнить оплату
		Payload:       payload,
		ProviderToken: "", // Для Stars должен быть пустой
		Currency:      string(currency),
		Prices: []tgbotapi.LabeledPrice{
			{
				Label:  common.MESSAGES[lang]["premium"] + " " + subscriptionPeriodLabel, // Если появятся другие планы, то нужно будет формировать динамически
				Amount: int(pricing.Price),
			},
		},
		StartParameter: "invoice_lock", // Чтобы предотвратить оплату из других чатов
		ProtectContent: true,           // Чтобы предотвратить оплату из других чатов
	}, nil
}

func SendTgInvoice(ctx context.Context, tgBotAPI *tgbotapi.Bot, tgInvoiceParams *tgbotapi.SendInvoiceParams) error {
	_, err := tgBotAPI.SendInvoice(ctx, tgInvoiceParams)

	return err
}

func getSubscriptionPeriodLabel(lang string, subscriptionPeriodUnit subPkg.SubscriptionPeriodUnit, subscriptionPeriodCount int64) string {
	if subscriptionPeriodUnit == subPkg.SubscriptionPeriodUnitLifetime {
		return common.MESSAGES[lang]["lifetime"]
	}
	return common.MESSAGES[lang]["for"] + " " + fmt.Sprint(subscriptionPeriodCount) + " " + common.MESSAGES[lang][string(subscriptionPeriodUnit)+"Short"]
}
