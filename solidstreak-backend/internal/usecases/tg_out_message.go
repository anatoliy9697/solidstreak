package usecases

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	invPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice"
	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
)

func NewErrorTgMessageParams(r common.Resources, tcTgID int64, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = common.GetDefaultLang()
	}

	return tu.Message(
		tu.ID(tcTgID),
		common.MESSAGES[lang]["smthWrongMsg"],
	)
}

func NewDefaultSuccessTgMessageParams(r common.Resources, tcTgID int64, userTgFirstName string, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = common.GetDefaultLang()
	}

	return tu.Message(
		tu.ID(tcTgID),
		fmt.Sprintf(common.MESSAGES[lang]["helloMsg"], userTgFirstName),
	).WithReplyMarkup(tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(common.MESSAGES[lang]["open"]).WithWebApp(tu.WebAppInfo(r.WebAppURL)),
		),
	))
}

func NewSubscriptionPurchaseSuccessTgMessageParams(r common.Resources, tcTgID int64, renew bool, subscriptionEvent *subPkg.SubscriptionEvent, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = common.GetDefaultLang()
	}

	var actionLabel string
	if renew {
		actionLabel = common.MESSAGES[lang]["renewed"]
	} else {
		actionLabel = common.MESSAGES[lang]["purchased"]
	}

	subscriptionPlanLabel := common.MESSAGES[lang][subscriptionEvent.SubscriptionPlanCode]

	subscriptionPeriodLabel := buildSubscriptionPeriodLabel(lang, subscriptionEvent.SubscriptionPeriodUnit, subscriptionEvent.SubscriptionPeriodCount)

	var subscriptionActiveUntilLabel string
	if subscriptionEvent.SubscriptionPeriodUnit != subPkg.SubscriptionPeriodUnitLifetime {
		subscriptionActiveUntilLabel = fmt.Sprintf(common.MESSAGES[lang]["subscriptionActiveUntil"], subscriptionEvent.SubscriptionPeriodFinishDate.DDMMYYYYString()) + " (UTC)\n\n"
	}

	return tu.Message(
		tu.ID(tcTgID),
		fmt.Sprintf(
			common.MESSAGES[lang]["subscriptionPurchasedMsg"],
			actionLabel,
			subscriptionPlanLabel,
			subscriptionPeriodLabel,
			subscriptionActiveUntilLabel,
		),
	)
}

func SendTgMessage(ctx context.Context, tgBotAPI *tgbotapi.Bot, params *tgbotapi.SendMessageParams) error {
	_, err := tgBotAPI.SendMessage(ctx, params)

	return err
}

func NewTgInvoiceParams(
	r common.Resources,
	tcTgID int64,
	lang string,
	invoice *invPkg.Invoice,
	subscriptionEvent *subPkg.SubscriptionEvent,
	invoiceExpiresIn time.Duration,
) (*tgbotapi.SendInvoiceParams, error) {
	if lang == "" {
		lang = common.GetDefaultLang()
	}

	payload := BuildInvoicePayloadString(invoice, subscriptionEvent)

	subscriptionPeriodLabel := buildSubscriptionPeriodLabel(lang, subscriptionEvent.SubscriptionPeriodUnit, subscriptionEvent.SubscriptionPeriodCount)
	return &tgbotapi.SendInvoiceParams{
		ChatID:        tgbotapi.ChatID{ID: tcTgID},
		Title:         common.MESSAGES[lang]["premiumSubscription"],                                                                                                                                  // Если появятся другие планы, то нужно будет формировать динамически
		Description:   common.MESSAGES[lang]["accessToPremium"] + " " + subscriptionPeriodLabel + ". " + fmt.Sprintf(common.MESSAGES[lang]["invoiceExpirationMsg"], int(invoiceExpiresIn.Minutes())), // Если появятся другие планы, то нужно будет формировать динамически
		Payload:       payload,
		ProviderToken: "", // Для Stars должен быть пустой
		Currency:      string(invoice.Currency),
		Prices: []tgbotapi.LabeledPrice{
			{
				Label:  common.MESSAGES[lang]["premium"] + " " + subscriptionPeriodLabel, // Если появятся другие планы, то нужно будет формировать динамически
				Amount: int(invoice.Amount),
			},
		},
		StartParameter: "invoice_lock", // Чтобы предотвратить оплату из других чатов
		ProtectContent: true,           // Чтобы предотвратить оплату из других чатов
	}, nil
}

func SendTgInvoice(ctx context.Context, tgBotAPI *tgbotapi.Bot, params *tgbotapi.SendInvoiceParams) (*tgbotapi.Message, error) {
	return tgBotAPI.SendInvoice(ctx, params)
}

func NewTgAnswerPreCheckoutQuery(pcqID string, ok bool, lang string) *tgbotapi.AnswerPreCheckoutQueryParams {
	if lang == "" {
		lang = common.GetDefaultLang()
	}

	params := &tgbotapi.AnswerPreCheckoutQueryParams{
		PreCheckoutQueryID: pcqID,
		Ok:                 ok,
	}

	if !ok {
		params.ErrorMessage = common.MESSAGES[lang]["smthWrongMsg"]
	}

	return params
}

func SendTgAnswerPreCheckoutQuery(ctx context.Context, tgBotAPI *tgbotapi.Bot, params *tgbotapi.AnswerPreCheckoutQueryParams) error {
	return tgBotAPI.AnswerPreCheckoutQuery(ctx, params)
}

func buildSubscriptionPeriodLabel(lang string, subscriptionPeriodUnit subPkg.SubscriptionPeriodUnit, subscriptionPeriodCount int64) string {
	if subscriptionPeriodUnit == subPkg.SubscriptionPeriodUnitLifetime {
		return common.MESSAGES[lang]["lifetime"]
	}
	return common.MESSAGES[lang]["for"] + " " + fmt.Sprint(subscriptionPeriodCount) + " " + common.MESSAGES[lang][string(subscriptionPeriodUnit)+"Short"]
}

func NewDeleteTgMessageParams(tcTgID int64, tgMessageID int) *tgbotapi.DeleteMessageParams {
	return &tgbotapi.DeleteMessageParams{
		ChatID:    tgbotapi.ChatID{ID: tcTgID},
		MessageID: tgMessageID,
	}
}

func DeleteTgMessage(ctx context.Context, tgBotAPI *tgbotapi.Bot, params *tgbotapi.DeleteMessageParams) error {
	return tgBotAPI.DeleteMessage(ctx, params)
}

func NewExpiredInvoiceTgMessageParams(r common.Resources, tcTgID int64, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = common.GetDefaultLang()
	}

	return tu.Message(
		tu.ID(tcTgID),
		common.MESSAGES[lang]["expiredInvoiceMsg"],
	)
}
