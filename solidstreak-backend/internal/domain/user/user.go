package user

import "time"

type User struct {
	ID          int64     `json:"id"`
	TgID        int64     `json:"tgId"`
	TgUsername  string    `json:"tgUsername"`
	TgFirstName string    `json:"tgFirstName"`
	TgLastName  string    `json:"tgLastName"`
	TgLangCode  string    `json:"tgLangCode"`
	LangCode    string    `json:"langCode"`
	TgIsBot     bool      `json:"tgIsBot"`
	CreatedAt   time.Time `json:"createdAt"`
}

func NewUser(tgID int64, tgUsername, tgFirstName, tgLastName, tgLangCode, langCode string, tgIsBot bool) *User {
	return &User{
		TgID:        tgID,
		TgUsername:  tgUsername,
		TgFirstName: tgFirstName,
		TgLastName:  tgLastName,
		TgLangCode:  tgLangCode,
		LangCode:    langCode,
		TgIsBot:     tgIsBot,
		CreatedAt:   time.Now(),
	}
}
