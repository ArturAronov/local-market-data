package main

import (
	"context"
	"log"
	company_info "market-data/src/company-info"
	"market-data/src/params"
	"market-data/src/user"
	"market-data/src/utils"

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
	marketDb := dbs[utils.COMPANY_INFO]

	defer userDb.Close()
	defer marketDb.Close()

	userRepo := user.NewRepository(userDb)
	companyRepo := company_info.NewRepository(marketDb)
	companyCtrl := company_info.NewController(companyRepo)

	params.ParamsResover(userRepo, companyCtrl)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

}
