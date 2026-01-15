package params

import (
	"log"
	company_info "market-data/src/company-info"
	initcmd "market-data/src/init"

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

func ParamsResover() {
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
		company_info.GetCompanyFactsC(1018724)
	case paramsMap[INIT]:
		runInitErr := initcmd.RunInit(args[2:])
		if runInitErr != nil {
			log.Fatalf("ParamsResover: init error %v\n%v\n", runInitErr, os.Stderr)
		}
	case paramsMap[UPDATE]:
		runInitErr := RunUpdate(args[2:])
		if runInitErr != nil {
			log.Fatalf("ParamsResover: init error %v\n%v\n", runInitErr, os.Stderr)
		}
	default:
		log.Fatalf("Unknown subcommand: %v\n", args[1])
	}
}
