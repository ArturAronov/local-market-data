package main

import (
	"context"
	"log"
	"market-data/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	dbs, dbsErr := utils.InitDB()
	if dbsErr != nil {
		log.Fatalln(dbsErr)
	}

	userDb := dbs[utils.USER_EMAIL]
	marketDb := dbs[utils.COMPANY_TICKERS]

	defer userDb.Close()
	defer marketDb.Close()

	utils.ParamsResover()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

}
