package company_info

import (
	"encoding/json"
	"fmt"
	"log"

	"market-data/src/user"
	"market-data/src/utils"
)

func GetCompanyTickersC(email string) int {
	url := "https://www.sec.gov/files/company_tickers.json"

	body, res, bodyErr := utils.HttpReq(email, url)
	if bodyErr != nil {
		log.Fatalf("[GetCompanyTickersC] failed to handle request %v", bodyErr)
	}

	var secRes SecEntryRes
	jsonErr := json.Unmarshal(body, &secRes)
	if jsonErr != nil {
		log.Fatalf("[GetCompanyTickersC] Failed to unmarshal JSON:\n%v\n", jsonErr)
	}

	InsertTickerInfoR(&secRes)

	return res.StatusCode
}

func GetCompanyFactsC(cik int) {
	cikStr := fmt.Sprintf("%010d", cik)
	url := fmt.Sprintf("https://data.sec.gov/api/xbrl/companyfacts/CIK%s.json", cikStr)

	email, emailErr := user.GetUserEmail()
	if emailErr != nil {
		log.Fatalf("[GetCompanyFactsC] Error getting user email: %v\n", emailErr)
	}

	body, _, bodyErr := utils.HttpReq(*email, url)
	if bodyErr != nil {
		log.Fatalf("[GetCompanyFactsC] failed to handle request %v", bodyErr)
	}

	var secRes CompanyFacts
	jsonErr := json.Unmarshal(body, &secRes)
	if jsonErr != nil {
		log.Fatalf("[GetCompanyFactsC] failed to unmarshal JSON:\n%v\n,", jsonErr)
	}

	dbCompany, dbCompnayErr := GetCompanyR(cik)
	if dbCompnayErr != nil {
		log.Fatal(dbCompnayErr)
	}

	if dbCompany == nil {
		EnterCompanyInfo(cik)
	}
}
