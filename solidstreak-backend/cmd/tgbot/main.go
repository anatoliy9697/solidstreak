package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	var err error

	err = godotenv.Load("./pkg/config/.env")
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
}
