package initcmd

import (
	"flag"
	"fmt"
	company_ticker "market-data/src/company-tickers"
)

func RunInit(args []string) error {
	var email string

	initFlags := flag.NewFlagSet("init", flag.ContinueOnError)
	initFlags.StringVar(&email, "e", "", "Alias of --email")
	initFlags.StringVar(&email, "email", "", "Your email for userAgent header")

	initFlagsErr := initFlags.Parse(args)
	if initFlagsErr != nil {
		return initFlagsErr
	}

	initFlags.Args()

	secResponseCode, _, secResponseErr := company_ticker.GetCompanyTickersC(email)

	fmt.Println(secResponseCode)

	if secResponseErr != nil {
		return secResponseErr
	}

	if secResponseCode == 403 {
		return fmt.Errorf("[RunInit] SEC response returned status code 403 (Forbidden), this is likely due to wrong email address provided")
	}

	if secResponseCode < 300 {
		// TODO:
		// Store email in local db
		// Add email to cache
	}

	return nil
}
