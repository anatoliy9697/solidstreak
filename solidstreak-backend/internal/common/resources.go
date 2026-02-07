package common

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	h "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/habit/repo"
	tc "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat/repo"
	usr "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user/repo"
)

type Resources struct {
	TgBotAPIToken string
	Logger        *slog.Logger
	TgBotAPI      *tgbotapi.BotAPI
	UsrRepo       usr.Repo
	TCRepo        tc.Repo
	HabitRepo     h.Repo
}
