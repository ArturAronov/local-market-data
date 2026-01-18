package company_info

import (
	"encoding/json"
	"fmt"
	"log"

	company_metadata "market-data/src/company-metadata"
	"market-data/src/user"
	"market-data/src/utils"
	"strings"
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
		var companyData Company
		submissions := company_metadata.GetSubmissionsC(cik)
		reportDates := company_metadata.GetLatestReportDate(submissions.Filings.Recent)

		companyData.Cik = cik
		companyData.Sic = submissions.SIKStr
		companyData.Ticker = secRes.EntityName
		companyData.Phone = submissions.Phone
		companyData.EntryType = submissions.EntryType
		companyData.OwnerOrg = submissions.OwnerOrg
		companyData.Exchanges = strings.Join(submissions.Exchanges, ",")
		companyData.Description = submissions.Description
		companyData.FiscalYearEnd = submissions.FiscalYearEnd

		if reportDates != nil {
			companyData.Latest10k = reportDates.Latest10k
			companyData.Latest10q = reportDates.Latest10q
		}

		updateErr := UpdateCompanyR(companyData)
		if updateErr != nil {
			log.Fatal(updateErr)
		}
	}
}
