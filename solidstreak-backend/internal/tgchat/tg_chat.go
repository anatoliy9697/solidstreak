package tgchat

import "time"

type Chat struct {
	TgID      int64     `json:"tgId"`
	UserID    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}
