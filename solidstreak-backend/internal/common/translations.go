package common

var MESSAGES = map[string]map[string]string{
	"en": {
		"smthWrong": "Something went wrong\nPlease try again later",
		"helloMsg":  "Hello, %s!\nPush \"Open\" button to start using bot",
		"open":      "Open",
	},
	"ru": {
		"smthWrong": "Что-то пошло не так\nПожалуйста, попробуйте позже",
		"helloMsg":  "Привет, %s!\nНажмите кнопку \"Открыть\", чтобы начать пользоваться ботом",
		"open":      "Открыть",
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
