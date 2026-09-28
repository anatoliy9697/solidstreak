package common

var MESSAGES = map[string]map[string]string{
	"en": {
		"smthWrongMsg":             "Something went wrong\nPlease try again later or contact support (@solidstreak_chat, @avasin_dev)",
		"helloMsg":                 "Hello, %s!\nPush \"Open\" button to start using bot",
		"open":                     "Open",
		"premiumSubscription":      "Premium subscription",
		"accessToPremium":          "Access to premium features",
		"for":                      "for",
		"monthShort":               "mo.",
		"yearShort":                "yr.",
		"lifetime":                 "for life",
		"premium":                  "Premium",
		"invoiceExpirationMsg":     "Invoice must be paid within %d min, otherwise it will be automatically canceled",
		"subscriptionPurchasedMsg": "Congratulations! You have successfully %s a %s subscription %s",
		"purchased":                "purchased",
		"renewed":                  "renewed",
		"subscriptionActiveUntil":  "Your subscription is active until %s",
		"expiredInvoiceMsg":        "Your invoice has expired\nPlease try again to purchase or renew your subscription",
	},
	"ru": {
		"smthWrongMsg":             "Что-то пошло не так\nПожалуйста, попробуйте позже или обратитесь в поддержку (@solidstreak_chat, @avasin_dev)",
		"helloMsg":                 "Привет, %s!\nНажмите кнопку \"Открыть\", чтобы начать пользоваться ботом",
		"open":                     "Открыть",
		"premiumSubscription":      "Premium-подписка",
		"accessToPremium":          "Доступ к premium-функциям",
		"for":                      "на",
		"monthShort":               "мес.",
		"yearShort":                "г.",
		"lifetime":                 "навсегда",
		"premium":                  "Premium",
		"invoiceExpirationMsg":     "Счет необходимо оплатить в течение %d мин, иначе он будет автоматически отменен",
		"subscriptionPurchasedMsg": "Поздравляем! Вы успешно %s %s-подписку %s",
		"purchased":                "приобрели",
		"renewed":                  "продлили",
		"subscriptionActiveUntil":  "Ваша подписка действует до %s",
		"expiredInvoiceMsg":        "Ваш счет на оплату просрочен\nПожалуйста, попробуйте снова приобрести или продлить подписку",
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
