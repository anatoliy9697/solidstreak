package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/control/tgbot"
	tcRepo "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat/repo"
	usrRepo "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user/repo"
)

func main() {
	var err error

	viper.SetConfigFile("./pkg/config/config.yaml")
	if err = viper.ReadInConfig(); err != nil {
		log.Fatal("error reading config file: " + err.Error())
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(viper.GetInt("log_level"))}))

	defer func() {
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	logger.Info("solid Streak initialization...")

	err = godotenv.Load("./pkg/config/.env")
	if err != nil {
		return
	}
	if os.Getenv("POSTGRES_CONN_STRING") == "" {
		err = fmt.Errorf("POSTGRES_CONN_STRING is not set")
		return
	}

	// Creating main context that will be cancelled on SIGINT or SIGTERM
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Storage connection pool initialization
	var pgPool *pgxpool.Pool
	if pgPool, err = pgxpool.New(mainCtx, os.Getenv("POSTGRES_CONN_STRING")); err != nil {
		return
	}
	defer pgPool.Close()

	var tgBotAPI *tgbotapi.BotAPI
	if tgBotAPI, err = tgbotapi.NewBotAPI(os.Getenv("TG_BOT_API_TOKEN")); err != nil {
		return
	}
	// tgBotAPI.Debug = true

	goroutinesDoneCh := make(chan struct{}, 1)

	go tgbot.EventFetcher{
		TgBotUpdsOffset:  viper.GetInt("tg_bot_upds_offset"),
		TgBotUpdsTimeout: viper.GetInt("tg_bot_upds_timeout"),
		TgBotAPI:         tgBotAPI,
		Logger:           logger,
		UsrRepo:          usrRepo.Init(mainCtx, pgPool),
		TcRepo:           tcRepo.Init(mainCtx, pgPool),
	}.Run(mainCtx, goroutinesDoneCh)

	logger.Info("solid Streak started")

	// Keeping alive
	<-mainCtx.Done()

	// Waiting for goroutines to finish
	<-goroutinesDoneCh

	logger.Info("solid Streak stopped")
}
