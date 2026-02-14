package params

import (
	"flag"
	"log"
	company_info "market-data/src/company-info"
	"strconv"
)

type CompanyInfoFlags struct {
	cik    int
	ticker string
}

var companyInfoFlags CompanyInfoFlags

func runCompanyInfo(args []string, email string, companyCtrl *company_info.Controller) error {
	initFlags := flag.NewFlagSet("update", flag.ContinueOnError)
	initFlags.IntVar(&companyInfoFlags.cik, "c", 0, "Alias of --cik")
	initFlags.IntVar(&companyInfoFlags.cik, "cik", 0, "Company CIK number")
	initFlags.StringVar(&companyInfoFlags.ticker, "t", "", "Alias for --ticker")
	initFlags.StringVar(&companyInfoFlags.ticker, "ticker", "", "Database containing company ticker and CIK data")

	initFlagsErr := initFlags.Parse(args)
	if initFlagsErr != nil {
		return initFlagsErr
	}
	initFlags.Args()

	if companyInfoFlags.ticker != "" {
		cik := companyCtrl.GetCikByTickerC(companyInfoFlags.ticker)
		cikInt, cikIntErr := strconv.Atoi(cik)
		if cikIntErr != nil {
			log.Fatal(cikIntErr)
		}

		companyInfoFlags.cik = cikInt
	}

	if companyInfoFlags.cik > 0 {
		companyCtrl.GetCompanyFinancialReportsC(companyInfoFlags.cik, email)
	}

	return nil
}
