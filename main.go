package main

import (
	"context"
	"database/sql"
	"log"
	user_agent "market-data/src/user-agent"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, dbErr := sql.Open("sqlite3", "./data/market-data.db")
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	defer db.Close()
	user_agent.InitDB(ctx, db)

	log.Println("Server starting on :7700")

}
