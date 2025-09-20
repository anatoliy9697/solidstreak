package tgbot

import (
	"context"

	"github.com/google/uuid"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type EventFetcher struct {
	TgBotUpdsOffset  int
	TgBotUpdsTimeout int
	MaxEventHandlers int
	Res              common.Resources
}

func (ef EventFetcher) Run(ctx context.Context, doneCh chan struct{}) {
	defer func() { doneCh <- struct{}{} }()

	ef.Res.Logger.Info("event fetcher started")

	updConfig := tgbotapi.NewUpdate(ef.TgBotUpdsOffset)
	updConfig.Timeout = ef.TgBotUpdsTimeout

	upds := ef.Res.TgBotAPI.GetUpdatesChan(updConfig)

	handlers := make(map[string]struct{}, ef.MaxEventHandlers)
	handlerDoneCh := make(chan string, ef.MaxEventHandlers)
	handlerCode := ""

loop:
	for {
		select {

		// Event fetcher stopping
		case <-ctx.Done():
			ef.Res.Logger.Info("event fetcher shutting down initiated")
			break loop

		// Event handler finished
		case handlerCode = <-handlerDoneCh:
			delete(handlers, handlerCode)
			ef.Res.Logger.Debug("event handler finished", "handlerCode", handlerCode)

		// New update received
		case upd := <-upds:
			// Ignore non-message and non-callback updates
			if upd.Message == nil && upd.CallbackQuery == nil {
				continue
			}
			ef.Res.Logger.Info("new update received", "update", upd)
			if len(handlers) >= ef.MaxEventHandlers {
				ef.Res.Logger.Warn("max event handlers limit reached, waiting for a handler to finish")
				handlerCode = <-handlerDoneCh
				delete(handlers, handlerCode)
			}
			handlerCode = uuid.NewString()[:8]
			handlers[handlerCode] = struct{}{}
			ef.Res.Logger.Debug("running new event handler", "handlerCode", handlerCode)
			go EventHandler{
				Code: handlerCode,
				Res:  ef.Res,
			}.Run(handlerDoneCh, &upd)
		}
	}

	if len(handlers) > 0 {
		ef.Res.Logger.Info("waiting for event handlers to finish")
	}
	for len(handlers) > 0 {
		handlerCode = <-handlerDoneCh
		delete(handlers, handlerCode)
	}

	ef.Res.Logger.Info("event fetcher stopped")
}
