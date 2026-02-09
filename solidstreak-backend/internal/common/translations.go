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

func ToLocalLang(langCode string) string {
	for _, lang := range LANGS {
		if langCode == lang {
			return lang
		}
	}
	return GetDefaultLang()
}
