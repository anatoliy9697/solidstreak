package common

const supportMsgEn = `If you have any questions or issues with the application, please describe your question or issue in the support chat (@solidstreak_chat) or directly to the developer (@avasin_dev)

If you have problems with subscription payment, use the <code>/paysupport</code> command`

const supportMsgRu = `Если у вас есть вопросы или проблемы в работе приложения, пожалуйста, опишите ваш вопрос или проблему в чате с поддержкой (@solidstreak_chat) или напрямую разработчику (@avasin_dev)

Если у вас проблемы с оплатой подписки, воспользуйтесь командой <code>/paysupport</code>`

const paySupportMsgEn = `To submit a payment support ticket, send the <code>/paysupport</code> command together with a description of the issue

Please describe the issue in detail. If possible, include:
• the date and time of the incident;
• the transaction ID;
• the name of the purchased service or subscription;
• any other information that may help us investigate the issue

<b>Example:</b>

<code>/paysupport My subscription was not activated after payment. Date and time: 2026-10-06 15:23:00+03:00. Transaction ID: stxGZuU...</code>

We will review your ticket within 3 calendar days and inform you of the outcome in this chat`

const paySupportMsgRu = `Для регистрации обращения по вопросу или проблеме, связанной с оплатой, отправьте команду <code>/paysupport</code> вместе с описанием проблемы

Постарайтесь описать проблему подробно. По возможности укажите:
• дату и время инцидента;
• идентификатор транзакции;
• название оплаченной услуги или подписки;
• другую информацию, которая может помочь разобраться в ситуации

<b>Пример:</b>

<code>/paysupport Не активировалась подписка после оплаты. Дата и время: 2026-10-06 15:23:00+03:00. Идентификатор транзакции: stxGZuU...</code>

Мы рассмотрим ваше обращение в течение 3 календарных дней и сообщим о результате в этом чате`

const paySupportTicketCreatedMsgEn = `Your payment support ticket has been created: <code>%s</code>

We will review your ticket within 3 calendar days and inform you of the outcome in this chat`

const paySupportTicketCreatedMsgRu = `Ваше обращение по вопросу оплаты создано: <code>%s</code>

Мы рассмотрим ваше обращение в течение 3 календарных дней и сообщим о результате в этом чате`

const alreadyOpenTicketMsgEn = `You already have an active support ticket: <code>%s</code>

Please wait for a response to it before creating a new ticket`

const alreadyOpenTicketMsgRu = `У вас уже есть активное обращение: <code>%s</code>

Дождитесь обратной связи по нему, прежде чем создавать новое обращение`

var MESSAGES = map[string]map[string]string{
	"en": {
		"smthWrongMsg":               "Something went wrong\nPlease try again later or contact support (@solidstreak_chat, @avasin_dev)",
		"helloMsg":                   "Hello, %s!\nPush \"Open\" button to start using bot",
		"open":                       "Open",
		"premiumSubscription":        "Premium subscription",
		"accessToPremium":            "Access to premium features",
		"for":                        "for",
		"monthShort":                 "mo.",
		"yearShort":                  "yr.",
		"lifetime":                   "for life",
		"premium":                    "Premium",
		"invoiceExpirationMsg":       "Invoice (%s) must be paid within %d min, otherwise it will be automatically canceled",
		"subscriptionPurchasedMsg":   "Congratulations! You have successfully %s a %s subscription %s",
		"purchased":                  "purchased",
		"renewed":                    "renewed",
		"subscriptionActiveUntil":    "Your subscription is active until %s",
		"expiredInvoiceMsg":          "Your invoicehas expired\nPlease try again to purchase or renew your subscription",
		"supportMsg":                 supportMsgEn,
		"paySupportMsg":              paySupportMsgEn,
		"paySupportTicketCreatedMsg": paySupportTicketCreatedMsgEn,
		"alreadyOpenTicketMsg":       alreadyOpenTicketMsgEn,
		"refundMsg":                  "Refund of %d %s has been processed for invoice <code>%s</code>",
	},
	"ru": {
		"smthWrongMsg":               "Что-то пошло не так\nПожалуйста, попробуйте позже или обратитесь в поддержку (@solidstreak_chat, @avasin_dev)",
		"helloMsg":                   "Привет, %s!\nНажмите кнопку \"Открыть\", чтобы начать пользоваться ботом",
		"open":                       "Открыть",
		"premiumSubscription":        "Premium-подписка",
		"accessToPremium":            "Доступ к premium-функциям",
		"for":                        "на",
		"monthShort":                 "мес.",
		"yearShort":                  "г.",
		"lifetime":                   "навсегда",
		"premium":                    "Premium",
		"invoiceExpirationMsg":       "Счет (%s) необходимо оплатить в течение %d мин, иначе он будет автоматически отменен",
		"subscriptionPurchasedMsg":   "Поздравляем! Вы успешно %s %s-подписку %s",
		"purchased":                  "приобрели",
		"renewed":                    "продлили",
		"subscriptionActiveUntil":    "Ваша подписка действует до %s",
		"expiredInvoiceMsg":          "Ваш счет на оплату просрочен\nПожалуйста, попробуйте снова приобрести или продлить подписку",
		"supportMsg":                 supportMsgRu,
		"paySupportMsg":              paySupportMsgRu,
		"paySupportTicketCreatedMsg": paySupportTicketCreatedMsgRu,
		"alreadyOpenTicketMsg":       alreadyOpenTicketMsgRu,
		"refundMsg":                  "Выполнен возврат в размере %d %s по счёту <code>%s</code>",
	},
}

var LANGS = []string{"en", "ru"}

func GetDefaultLang() string {
	return LANGS[0]
}

func IsSupportedLang(langCode string) bool {
	for _, lang := range LANGS {
		if langCode == lang {
			return true
		}
	}
	return false
}

func ToLocalLang(langCode string) string {
	if IsSupportedLang(langCode) {
		return langCode
	} else {
		return GetDefaultLang()
	}
}
