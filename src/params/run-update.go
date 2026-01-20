package params

import (
	"flag"
	"fmt"

	company_info "market-data/src/company-info"
	"market-data/src/user"
)

type UserFlagsT struct {
	email          string
	companyTickers bool
}

var userFlags UserFlagsT

func RunUpdate(args []string, userRepo *user.Repository, companyCtrl *company_info.Controller) error {
	initFlags := flag.NewFlagSet("update", flag.ContinueOnError)
	initFlags.StringVar(&userFlags.email, "e", "", "Alias of --email")
	initFlags.StringVar(&userFlags.email, "email", "", "Your email for userAgent header")
	initFlags.BoolVar(&userFlags.companyTickers, "ct", false, "Alias for --company-tickers")
	initFlags.BoolVar(&userFlags.companyTickers, "company-tickers", false, "Database containing company ticker and CIK data")

	initFlagsErr := initFlags.Parse(args)
	if initFlagsErr != nil {
		return initFlagsErr
	}

	initFlags.Args()

	switch {
	case userFlags.email != "":
		userRepo.InsertUserEmail(userFlags.email)
		fallthrough
	case userFlags.companyTickers:
		email, emailErr := userRepo.GetUserEmail()
		if emailErr != nil {
			return emailErr
		}

		companyCtrl.GetCompanyTickersC(*email)
		// fallthrough
	default:
		return fmt.Errorf("Flags can't be empty")
	}

	return nil
}
