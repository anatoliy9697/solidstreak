package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common/resources"
	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
)

func SendReplyMsg(r resources.Resources, tc *tcPkg.Chat, msgText string) error {
	msg := tgbotapi.NewMessage(tc.TgID, msgText)

	_, err := r.TgBotAPI.Send(msg)

	return err
}
