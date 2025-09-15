package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/tgchat"
	tcRepo "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/tgchat/repo"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/user"
	usrRepo "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/user/repo"
)

func main() {
	var err error

	err = godotenv.Load("./solidstreak-backend/pkg/config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("POSTGRES_CONN_STRING") == "" {
		log.Fatal("POSTGRES_CONN_STRING is not set")
	}

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var pgPool *pgxpool.Pool
	if pgPool, err = pgxpool.New(mainCtx, os.Getenv("POSTGRES_CONN_STRING")); err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer pgPool.Close()

	userRepo := usrRepo.Init(mainCtx, pgPool)
	tcRepo := tcRepo.Init(mainCtx, pgPool)

	u := &usrPkg.User{
		TgID:        123456789,
		TgUsername:  "testuser",
		TgFirstName: "Test",
		TgLastName:  "User",
		TgLangCode:  "en",
		TgIsBot:     false,
	}

	isUserExists, err := userRepo.IsExists(u)
	if err != nil {
		log.Fatal(err)
	}

	if isUserExists {
		err = userRepo.UpdateByTgId(u)
	} else {
		err = userRepo.Create(u)
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("User: %+v\n", u)

	tc := &tcPkg.Chat{
		TgID:      987654321,
		UserID:    u.ID,
		CreatedAt: time.Now(),
	}

	isTgChatExists, err := tcRepo.IsExistsByUserID(u.ID)
	if err != nil {
		log.Fatal(err)
	}

	if isTgChatExists {
		err = tcRepo.UpdateByUserID(tc)
	} else {
		err = tcRepo.Create(tc)
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("TgChat: %+v\n", tc)

	stop()
	log.Println("Shutting down gracefully...")
}
