package usecases

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common/i18n"

	invPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice"
	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
)

func NewErrorTgMessageParams(r common.Resources, tcTgID int64, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = i18n.GetDefaultLang()
	}

	return tu.Message(
		tu.ID(tcTgID),
		i18n.MESSAGES[lang]["smthWrongMsg"],
	)
}

func NewDefaultSuccessTgMessageParams(r common.Resources, tcTgID int64, userTgFirstName string, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = i18n.GetDefaultLang()
	}

	return tu.Message(
		tu.ID(tcTgID),
		fmt.Sprintf(i18n.MESSAGES[lang]["helloMsg"], userTgFirstName),
	).WithReplyMarkup(tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(i18n.MESSAGES[lang]["open"]).WithWebApp(tu.WebAppInfo(r.WebAppURL)),
		),
	))
}

func NewSubscriptionPurchaseSuccessTgMessageParams(r common.Resources, tcTgID int64, renew bool, subscriptionEvent *subPkg.SubscriptionEvent, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = i18n.GetDefaultLang()
	}

	var actionLabel string
	if renew {
		actionLabel = i18n.MESSAGES[lang]["renewed"]
	} else {
		actionLabel = i18n.MESSAGES[lang]["purchased"]
	}

	subscriptionPlanLabel := i18n.MESSAGES[lang][subscriptionEvent.SubscriptionPlanCode]

	subscriptionPeriodLabel := buildSubscriptionPeriodLabel(lang, subscriptionEvent.SubscriptionPeriodUnit, subscriptionEvent.SubscriptionPeriodCount)

	var subscriptionActiveUntilLabel string
	if subscriptionEvent.SubscriptionPeriodUnit != subPkg.SubscriptionPeriodUnitLifetime {
		subscriptionActiveUntilLabel = "\n\n" + fmt.Sprintf(i18n.MESSAGES[lang]["subscriptionActiveUntil"], subscriptionEvent.SubscriptionPeriodFinishDate.DDMMYYYYString()) + " (UTC)\n\n"
	}

	return tu.Message(
		tu.ID(tcTgID),
		fmt.Sprintf(
			i18n.MESSAGES[lang]["subscriptionPurchasedMsg"]+"%s",
			actionLabel,
			subscriptionPlanLabel,
			subscriptionPeriodLabel,
			subscriptionActiveUntilLabel,
		),
	)
}

func SendTgMessage(ctx context.Context, tgBotAPI *tgbotapi.Bot, params *tgbotapi.SendMessageParams) (*tgbotapi.Message, error) {
	msg, err := tgBotAPI.SendMessage(ctx, params)

	return msg, err
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
		lang = i18n.GetDefaultLang()
	}

	payload := BuildInvoicePayloadString(invoice, subscriptionEvent)

	subscriptionPeriodLabel := buildSubscriptionPeriodLabel(lang, subscriptionEvent.SubscriptionPeriodUnit, subscriptionEvent.SubscriptionPeriodCount)
	return &tgbotapi.SendInvoiceParams{
		ChatID:        tgbotapi.ChatID{ID: tcTgID},
		Title:         i18n.MESSAGES[lang]["premiumSubscription"],                                                                                                                                              // Если появятся другие планы, то нужно будет формировать динамически
		Description:   i18n.MESSAGES[lang]["accessToPremium"] + " " + subscriptionPeriodLabel + ". " + fmt.Sprintf(i18n.MESSAGES[lang]["invoiceExpirationMsg"], invoice.UUID, int(invoiceExpiresIn.Minutes())), // Если появятся другие планы, то нужно будет формировать динамически
		Payload:       payload,
		ProviderToken: "", // Для Stars должен быть пустой
		Currency:      string(invoice.Currency),
		Prices: []tgbotapi.LabeledPrice{
			{
				Label:  i18n.MESSAGES[lang]["premium"] + " " + subscriptionPeriodLabel, // Если появятся другие планы, то нужно будет формировать динамически
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
		lang = i18n.GetDefaultLang()
	}

	params := &tgbotapi.AnswerPreCheckoutQueryParams{
		PreCheckoutQueryID: pcqID,
		Ok:                 ok,
	}

	if !ok {
		params.ErrorMessage = i18n.MESSAGES[lang]["smthWrongMsg"]
	}

	return params
}

func SendTgAnswerPreCheckoutQuery(ctx context.Context, tgBotAPI *tgbotapi.Bot, params *tgbotapi.AnswerPreCheckoutQueryParams) error {
	return tgBotAPI.AnswerPreCheckoutQuery(ctx, params)
}

func buildSubscriptionPeriodLabel(lang string, subscriptionPeriodUnit subPkg.SubscriptionPeriodUnit, subscriptionPeriodCount int64) string {
	if subscriptionPeriodUnit == subPkg.SubscriptionPeriodUnitLifetime {
		return i18n.MESSAGES[lang]["lifetime"]
	}
	return i18n.MESSAGES[lang]["for"] + " " + fmt.Sprint(subscriptionPeriodCount) + " " + i18n.MESSAGES[lang][string(subscriptionPeriodUnit)+"Short"]
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

func NewExpiredInvoiceTgMessageParams(tcTgID int64, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = i18n.GetDefaultLang()
	}

	return tu.Message(
		tu.ID(tcTgID),
		i18n.MESSAGES[lang]["expiredInvoiceMsg"],
	)
}

func NewSupportTgMessageParams(tcTgID int64, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = i18n.GetDefaultLang()
	}

	msgParams := tu.Message(
		tu.ID(tcTgID),
		i18n.MESSAGES[lang]["supportMsg"],
	)
	msgParams.ParseMode = "HTML"

	return msgParams
}

func NewPaySupportTgMessageParams(tcTgID int64, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = i18n.GetDefaultLang()
	}

	msgParams := tu.Message(
		tu.ID(tcTgID),
		i18n.MESSAGES[lang]["paySupportMsg"],
	)
	msgParams.ParseMode = "HTML"

	return msgParams
}

func NewTgRefundStarPaymentParams(userTgID int64, tgPaymentChargeID string) *tgbotapi.RefundStarPaymentParams {
	return &tgbotapi.RefundStarPaymentParams{
		UserID:                  userTgID,
		TelegramPaymentChargeID: tgPaymentChargeID,
	}
}

func RefundStarPayment(ctx context.Context, tgBotAPI *tgbotapi.Bot, params *tgbotapi.RefundStarPaymentParams) error {
	return tgBotAPI.RefundStarPayment(ctx, params)
}

func NewPaySupportTicketCreatedTgMessageParams(tcTgID int64, ticketUUID string, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = i18n.GetDefaultLang()
	}

	msgParams := tu.Message(
		tu.ID(tcTgID),
		fmt.Sprintf(i18n.MESSAGES[lang]["paySupportTicketCreatedMsg"], ticketUUID),
	)
	msgParams.ParseMode = "HTML"

	return msgParams
}

func NewAlreadyOpenTicketTgMessageParams(tcTgID int64, ticketUUID string, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = i18n.GetDefaultLang()
	}

	msgParams := tu.Message(
		tu.ID(tcTgID),
		fmt.Sprintf(i18n.MESSAGES[lang]["alreadyOpenTicketMsg"], ticketUUID),
	)
	msgParams.ParseMode = "HTML"

	return msgParams
}

func NewRefundTgMessageParams(tcTgID int64, invoice *invPkg.Invoice, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = i18n.GetDefaultLang()
	}

	msgParams := tu.Message(
		tu.ID(tcTgID),
		fmt.Sprintf(i18n.MESSAGES[lang]["refundMsg"], invoice.Amount, invoice.Currency, invoice.UUID),
	)
	msgParams.ParseMode = "HTML"

	return msgParams
}

func NewTermsTgMessageParams(r common.Resources, tcTgID int64, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = i18n.GetDefaultLang()
	}

	msgParams := tu.Message(
		tu.ID(tcTgID),
		i18n.MESSAGES[lang]["termsCmdMsg"],
	).WithReplyMarkup(tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(i18n.MESSAGES[lang]["openTerms"]).WithWebApp(tu.WebAppInfo(r.WebAppURL + r.TermsPath + "/" + lang)),
		),
	))

	return msgParams
}

func NewSupportTicketAnswerTgMessageParams(tcTgID int64, ticketUUID string, ticketAnswer string, ticketStatusStr string, lang string) *tgbotapi.SendMessageParams {
	if lang == "" {
		lang = i18n.GetDefaultLang()
	}

	if ticketStatusStr == "closed" {
		ticketStatusStr = i18n.MESSAGES[lang]["supportTicketStatusClosed"]
	}

	msgParams := tu.Message(
		tu.ID(tcTgID),
		fmt.Sprintf(i18n.MESSAGES[lang]["supportTicketAnswerMsg"], ticketUUID, ticketAnswer, ticketStatusStr),
	)
	msgParams.ParseMode = "HTML"

	return msgParams
}
