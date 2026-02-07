package common

var MESSAGES = map[string]map[string]string{
	"en": {
		"smthWrong": "Something went wrong\nPlease try again later",
		"helloMsg":  "Hello, %s!\nPush \"Open\" button to start using bot",
	},
	"ru": {
		"smthWrong": "Что-то пошло не так\nПожалуйста, попробуйте позже",
		"helloMsg":  "Привет, %s!\nНажмите кнопку \"Open\", чтобы начать пользоваться ботом",
	},
}

var LANGS = []string{"en", "ru"}

func ToLocalLang(langCode string) string {
	for _, lang := range LANGS {
		if langCode == lang {
			return lang
		}
	}
	return "en"
}
