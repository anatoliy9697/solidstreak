package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	tgbotapi "github.com/mymmrac/telego"
	"github.com/spf13/viper"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/control/http"
	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/control/scheduler"
	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/control/tgbot"
	hRepo "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/habit/repo"
	invRepo "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice/repo"
	st "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/schedulertask"
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

	// creating main context that will be cancelled on SIGINT or SIGTERM
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// storage connection pool initialization
	var pgPool *pgxpool.Pool
	if pgPool, err = pgxpool.New(mainCtx, os.Getenv("POSTGRES_CONN_STRING")); err != nil {
		return
	}
	defer pgPool.Close()

	var tgBotAPI *tgbotapi.Bot
	if tgBotAPI, err = tgbotapi.NewBot(os.Getenv("TG_BOT_API_TOKEN")); err != nil {
		return
	}

	invRepoInstance := invRepo.Init(mainCtx, pgPool)
	resources := common.Resources{
		WebAppURL: os.Getenv("WEB_APP_URL"),
		Logger:    logger,
		TgBotAPI:  tgBotAPI,
		UsrRepo:   usrRepo.Init(mainCtx, pgPool),
		SubRepo: subRepo.Init(mainCtx, pgPool, subPkg.GetSubscriptionPlans(
			viper.GetInt64("basic_subscription.price.tg_stars_per_month"),
			viper.GetInt64("basic_subscription.price.tg_stars_per_year"),
			viper.GetInt64("basic_subscription.price.tg_stars_lifetime"),
			viper.GetInt64("basic_subscription.active_habits_limit"),
			viper.GetBool("basic_subscription.show_ads"),
			viper.GetInt64("premium_subscription.price.tg_stars_per_month"),
			viper.GetInt64("premium_subscription.price.tg_stars_per_year"),
			viper.GetInt64("premium_subscription.price.tg_stars_lifetime"),
			viper.GetInt64("premium_subscription.active_habits_limit"),
			viper.GetBool("premium_subscription.show_ads"),
		)),
		TCRepo:    tcRepo.Init(mainCtx, pgPool),
		HabitRepo: hRepo.Init(mainCtx, pgPool),
		InvRepo:   invRepoInstance,
	}

	goroutineDoneCh := make(chan struct{}, 2)

	// running event fetcher
	go tgbot.EventFetcher{
		MaxEventHandlers: viper.GetInt("max_tg_event_handlers"),
		Res:              resources,
	}.Run(mainCtx, goroutineDoneCh)

	// running scheduler
	go scheduler.Scheduler{
		Res:                    resources,
		TaskSources:            []st.TaskSource{invRepoInstance},
		MaxTaskHandlers:        viper.GetInt("scheduler.max_task_handlers"),
		TaskBatchSizePerSource: viper.GetInt("scheduler.task_batch_size_per_source"),
		TaskWaitingDuration:    time.Duration(viper.GetInt("scheduler.task_waiting_duration_ms")) * time.Millisecond,
		LockDuration:           time.Duration(viper.GetInt("scheduler.lock_duration_ms")) * time.Millisecond,
	}.Run(mainCtx, goroutineDoneCh)

	// running web server
	webServer := http.Server{
		Env:              os.Getenv("ENV"),
		CertFilePath:     os.Getenv("CERT_FILE_PATH"),
		KeyFilePath:      os.Getenv("KEY_FILE_PATH"),
		Addr:             os.Getenv("SERVER_ADDR"),
		InvoiceExpiresIn: time.Duration(viper.GetInt("invoice_expires_in_min")) * time.Minute,
		Res:              resources,
	}
	go webServer.Run(mainCtx, goroutineDoneCh)

	logger.Info("solid streak started")

	// keeping alive
	<-mainCtx.Done()

	// waiting for goroutines to finish
	<-goroutineDoneCh
	<-goroutineDoneCh

	logger.Info("solid streak stopped")
}

func validateConfigParams() error {
	if viper.GetInt("max_tg_event_handlers") <= 0 {
		return errors.New("max_tg_event_handlers should be greater than 0")
	}
	if viper.GetInt("invoice_expires_in_min") <= 0 {
		return errors.New("invoice_expires_in_min should be greater than 0")
	}
	if viper.GetInt("basic_subscription.active_habits_limit") <= 0 {
		return errors.New("basic_subscription.active_habits_limit should be greater than 0")
	}
	if viper.GetInt64("premium_subscription.price.tg_stars_per_month") <= 0 {
		return errors.New("premium_subscription.price.tg_stars_per_month should be greater than 0")
	}
	if viper.GetInt64("premium_subscription.price.tg_stars_per_year") <= 0 {
		return errors.New("premium_subscription.price.tg_stars_per_year should be greater than 0")
	}
	if viper.GetInt64("premium_subscription.price.tg_stars_lifetime") <= 0 {
		return errors.New("premium_subscription.price.tg_stars_lifetime should be greater than 0")
	}
	if viper.GetInt64("premium_subscription.active_habits_limit") <= 0 {
		return errors.New("premium_subscription.active_habits_limit should be greater than 0")
	}
	if viper.GetInt("scheduler.max_task_handlers") <= 0 {
		return errors.New("scheduler.max_task_handlers should be greater than 0")
	}
	if viper.GetInt("scheduler.task_batch_size_per_source") <= 0 {
		return errors.New("scheduler.task_batch_size_per_source should be greater than 0")
	}
	if viper.GetInt("scheduler.task_waiting_duration_ms") <= 0 {
		return errors.New("scheduler.task_waiting_duration_ms should be greater than 0")
	}
	if viper.GetInt("scheduler.lock_duration_ms") <= 0 {
		return errors.New("scheduler.lock_duration_ms should be greater than 0")
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
