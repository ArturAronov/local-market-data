package params

import (
	"log"
	company_info "market-data/src/company-info"
	initcmd "market-data/src/init"
	"market-data/src/user"

	"os"
)

type Params string

const (
	TEST   Params = "test"
	INIT   Params = "init"
	UPDATE Params = "update"
)

var paramsMap = map[Params]string{
	TEST:   string(TEST),
	INIT:   string(INIT),
	UPDATE: string(UPDATE),
}

func ParamsResover(userRepo *user.Repository, companyCtrl *company_info.Controller) {
	args := os.Args

	if len(args) == 1 {
		log.Println("Subcommand is expected. Use any of subcommands below:")

		for key := range paramsMap {
			log.Printf("- %s", paramsMap[key])
		}

		os.Exit(1)
	}

	if len(args) < 2 {
		log.Fatalf("Expected subcommand.")
	}

	switch args[1] {
	case paramsMap[TEST]:
		// submission_data.GetSubmissionDataC(1018724)
		email, err := userRepo.GetUserEmail()
		if err != nil {
			log.Fatalf("Failed to get user email for test: %v", err)
		}
		companyCtrl.GetCompanyFactsC(1018724, *email)
	case paramsMap[INIT]:
		runInitErr := initcmd.RunInit(args[2:], userRepo, companyCtrl)
		if runInitErr != nil {
			log.Fatalf("ParamsResover: init error %v\n%v\n", runInitErr, os.Stderr)
		}
	case paramsMap[UPDATE]:
		runInitErr := RunUpdate(args[2:], userRepo, companyCtrl)
		if runInitErr != nil {
			log.Fatalf("ParamsResover: init error %v\n%v\n", runInitErr, os.Stderr)
		}
	default:
		log.Fatalf("Unknown subcommand: %v\n", args[1])
	}
}
