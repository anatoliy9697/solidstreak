package main

import (
	"context"
	"errors"
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
	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
	subRepo "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription/repo"
	tcRepo "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat/repo"
	usrRepo "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user/repo"
)

func main() {
	var err error

	viper.SetConfigFile("./pkg/config/config.yaml")
	if err = viper.ReadInConfig(); err != nil {
		log.Fatal("error reading config file: " + err.Error())
	}
	if err = validateConfigParams(); err != nil {
		log.Fatal("invalid config parameter: " + err.Error())
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
	if err = validateEnvParams(); err != nil {
		log.Fatal("invalid env parameter: " + err.Error())
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
		SubRepo: subRepo.Init(mainCtx, pgPool, subPkg.GetSubscriptionPlans(
			viper.GetInt64("basic_sub_plan_active_habits_limit"),
			viper.GetInt64("premium_sub_plan_active_habits_limit"),
			viper.GetInt64("premium_sub_plan_price_stars_per_month"),
			viper.GetInt64("premium_sub_plan_price_stars_per_year"),
			viper.GetInt64("premium_sub_plan_price_stars_forever"),
		)),
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

func validateConfigParams() error {
	if viper.GetInt("max_event_handlers") <= 0 {
		return errors.New("max_event_handlers should be greater than 0")
	}
	if viper.GetInt("basic_sub_plan_active_habits_limit") <= 0 {
		return errors.New("basic_sub_plan_active_habits_limit should be greater than 0")
	}
	if viper.GetInt("premium_sub_plan_active_habits_limit") <= 0 {
		return errors.New("premium_sub_plan_active_habits_limit should be greater than 0")
	}
	if viper.GetInt("premium_sub_plan_price_stars_per_month") <= 0 {
		return errors.New("premium_sub_plan_price_stars_per_month should be greater than 0")
	}
	if viper.GetInt("premium_sub_plan_price_stars_per_year") <= 0 {
		return errors.New("premium_sub_plan_price_stars_per_year should be greater than 0")
	}
	if viper.GetInt("premium_sub_plan_price_stars_forever") <= 0 {
		return errors.New("premium_sub_plan_price_stars_forever should be greater than 0")
	}

	return nil
}

func validateEnvParams() error {
	if os.Getenv("POSTGRES_CONN_STRING") == "" {
		return errors.New("POSTGRES_CONN_STRING env variable is not set")
	}
	if os.Getenv("TG_BOT_API_TOKEN") == "" {
		return errors.New("TG_BOT_API_TOKEN env variable is not set")
	}
	if os.Getenv("WEB_APP_URL") == "" {
		return errors.New("WEB_APP_URL env variable is not set")
	}
	if os.Getenv("ENV") != "prod" && os.Getenv("ENV") != "dev" {
		return errors.New("ENV env variable should be either 'prod' or 'dev'")
	}
	if os.Getenv("SERVER_ADDR") == "" {
		return errors.New("SERVER_ADDR env variable is not set")
	}

	return nil
}
