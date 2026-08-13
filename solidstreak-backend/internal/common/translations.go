package common

var MESSAGES = map[string]map[string]string{
	"en": {
		"smthWrong":           "Something went wrong\nPlease try again later",
		"helloMsg":            "Hello, %s!\nPush \"Open\" button to start using bot",
		"open":                "Open",
		"premiumSubscription": "Premium subscription",
		"accessToPremium":     "Access to premium features",
		"month":               "for 1 month",
		"year":                "for 1 year",
		"lifetime":            "for life",
		"premium":             "Premium",
	},
	"ru": {
		"smthWrong":           "Что-то пошло не так\nПожалуйста, попробуйте позже",
		"helloMsg":            "Привет, %s!\nНажмите кнопку \"Открыть\", чтобы начать пользоваться ботом",
		"open":                "Открыть",
		"premiumSubscription": "Премиум подписка",
		"accessToPremium":     "Доступ к премиум функциям",
		"month":               "на 1 месяц",
		"year":                "на 1 год",
		"lifetime":            "навсегда",
		"premium":             "Премиум",
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
