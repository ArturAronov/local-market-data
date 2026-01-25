package params

import (
	"log"
	company_info "market-data/src/company-info"
	"market-data/src/user"
	"strings"

	"os"
)

type Params string

const (
	INIT         Params = "init"
	UPDATE       Params = "update"
	COMPANY_INFO Params = "ci"
)

var paramsMap = map[Params]string{
	INIT:         string(INIT),
	UPDATE:       string(UPDATE),
	COMPANY_INFO: string(COMPANY_INFO),
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

	if len(args) > 4 {
		log.Fatalf(
			"[ParamsResover] too many argumetns provided: %s",
			strings.Join(args[4:], ", "),
		)
	}

	switch args[1] {
	case paramsMap[INIT]:
		runInitErr := runInit(args[2:], userRepo, companyCtrl)
		if runInitErr != nil {
			log.Fatalf(
				"[ParamsResover] init error %v\n%v\n",
				runInitErr,
				os.Stderr,
			)
		}

	case paramsMap[UPDATE]:
		runUpdateErr := runUpdate(args[2:], userRepo, companyCtrl)
		if runUpdateErr != nil {
			log.Fatalf(
				"[ParamsResover] update error %v\n%v\n",
				runUpdateErr,
				os.Stderr,
			)
		}

	case paramsMap[COMPANY_INFO]:
		email, err := userRepo.GetUserEmail()
		if err != nil {
			log.Fatalf("[ParamsResover] Failed to get user email for: %v", err)
		}
		runCompanyInfoErr := runCompanyInfo(args[2:], *email, companyCtrl)
		if runCompanyInfoErr != nil {
			log.Fatalf(
				"[ParamsResover] update error %v\n%v\n",
				runCompanyInfoErr,
				os.Stderr,
			)
		}

	default:
		log.Fatalf("[ParamsResover] unknown subcommand: %v\n", args[1])
	}
}
