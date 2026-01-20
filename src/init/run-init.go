package initcmd

import (
	"flag"
	"fmt"

	company_info "market-data/src/company-info"
	"market-data/src/user"
)

func RunInit(args []string, userRepo *user.Repository, companyCtrl *company_info.Controller) error {
	var email string

	initFlags := flag.NewFlagSet("init", flag.ContinueOnError)
	initFlags.StringVar(&email, "e", "", "Alias of --email")
	initFlags.StringVar(&email, "email", "", "Your email for userAgent header")

	initFlagsErr := initFlags.Parse(args)
	if initFlagsErr != nil {
		return initFlagsErr
	}

	initFlags.Args()

	secResponseCode := companyCtrl.GetCompanyTickersC(email)

	if secResponseCode == 403 {
		return fmt.Errorf("[RunInit] SEC response returned status code 403 (Forbidden), this is likely due to wrong email address provided")
	}

	if secResponseCode < 300 {
		userRepo.InsertUserEmail(email)
	}

	return nil
}
