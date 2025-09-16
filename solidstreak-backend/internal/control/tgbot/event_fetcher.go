package tgbot

import (
	"context"
	"log/slog"

	tc "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat/repo"
	usr "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type EventFetcher struct {
	TgBotUpdsOffset  int
	TgBotUpdsTimeout int
	TgBotAPI         *tgbotapi.BotAPI
	Logger           *slog.Logger
	UsrRepo          usr.Repo
	TcRepo           tc.Repo
}

func (ef EventFetcher) Run(ctx context.Context, done chan struct{}) {
	defer func() { done <- struct{}{} }()

	ef.Logger.Info("event fetcher started")

	updConfig := tgbotapi.NewUpdate(ef.TgBotUpdsOffset)
	updConfig.Timeout = ef.TgBotUpdsTimeout

	upds := ef.TgBotAPI.GetUpdatesChan(updConfig)

loop:
	for {
		select {

		// Event fetcher stopping
		case <-ctx.Done():
			break loop

		// New update received
		case upd := <-upds:
			if upd.Message == nil {
				continue
			}
			ef.Logger.Info("new update received",
				slog.String("username", upd.Message.From.UserName),
				slog.Int64("chat_id", upd.Message.Chat.ID),
				slog.String("message_text", upd.Message.Text),
			)

		}
	}

	ef.Logger.Info("event fetcher stopped")
}
