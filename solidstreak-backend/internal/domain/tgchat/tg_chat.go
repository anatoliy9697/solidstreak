package tgchat

import "time"

type Chat struct {
	TgID      int64     `json:"tgId"`
	UserID    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewChat(tgID int64, userID int64) *Chat {
	return &Chat{
		TgID:      tgID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
}
