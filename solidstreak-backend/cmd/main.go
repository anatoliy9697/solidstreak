package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	tgbotapi "github.com/mymmrac/telego"
	"github.com/spf13/viper"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/control/http"
	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/control/tgbot"
	hRepo "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/habit/repo"
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
	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(viper.GetInt("log_level"))}))

	defer func() {
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	logger.Info("solid streak initialization...")

	err = godotenv.Load("./pkg/config/.env")
	if err != nil {
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

	var tgBotAPI *tgbotapi.Bot
	if tgBotAPI, err = tgbotapi.NewBot(os.Getenv("TG_BOT_API_TOKEN")); err != nil {
		return
	}

	resources := common.Resources{
		WebAppURL: os.Getenv("WEB_APP_URL"),
		Logger:    logger,
		TgBotAPI:  tgBotAPI,
		UsrRepo:   usrRepo.Init(mainCtx, pgPool),
		TCRepo:    tcRepo.Init(mainCtx, pgPool),
		HabitRepo: hRepo.Init(mainCtx, pgPool),
	}

	goroutineDoneCh := make(chan struct{}, 2)

	// Running event fetcher
	go tgbot.EventFetcher{
		MaxEventHandlers: viper.GetInt("max_event_handlers"),
		Res:              resources,
	}.Run(mainCtx, goroutineDoneCh)

	// Running web server
	webServer := http.Server{
		Env:          os.Getenv("ENV"),
		CertFilePath: os.Getenv("CERT_FILE_PATH"),
		KeyFilePath:  os.Getenv("KEY_FILE_PATH"),
		Addr:         os.Getenv("SERVER_ADDR"),
		Res:          resources,
	}
	go webServer.Run(mainCtx, goroutineDoneCh)

	logger.Info("solid streak started")

	// Keeping alive
	<-mainCtx.Done()

	// Waiting for goroutines to finish
	<-goroutineDoneCh
	<-goroutineDoneCh

	logger.Info("solid streak stopped")
}
