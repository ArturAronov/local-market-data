package company_info

import (
	"encoding/json"
	"fmt"
	"log"

	company_metadata "market-data/src/company-metadata"
	"market-data/src/utils"
)

type Controller struct {
	repo *Repository
}

func NewController(repo *Repository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) GetCompanyTickersC(email string) int {
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

	c.repo.InsertTickerInfoR(&secRes)

	return res.StatusCode
}

func (c *Controller) GetCompanyFinancialReportsC(cik int, email string) {
	cikStr := fmt.Sprintf("%010d", cik)
	url := fmt.Sprintf("https://data.sec.gov/api/xbrl/companyfacts/CIK%s.json", cikStr)

	body, _, bodyErr := utils.HttpReq(email, url)
	if bodyErr != nil {
		log.Fatalf("[GetCompanyFinancialReportsC] failed to handle request %v", bodyErr)
	}

	var secRes CompanyFacts
	jsonErr := json.Unmarshal(body, &secRes)
	if jsonErr != nil {
		log.Fatalf("[GetCompanyFinancialReportsC] failed to unmarshal JSON:\n%v\n,", jsonErr)
	}

	reportsCount, reportsCountErr := c.repo.CountCompanyReportDataR(cik)
	if reportsCountErr != nil {
		log.Fatal(reportsCountErr)
	}

	if reportsCount == 0 {
		c.EnterCompanyInfo(cik, email)
		c.EnterFinancialReports(secRes, email)

		return
	}

	company, companyErr := c.repo.GetCompanyR(cik)
	submissionData := company_metadata.GetSubmissionsC(cik, email)
	reportDates := company_metadata.GetLatestReportDate(submissionData.Filings.Recent)

	if *company.Latest10k != reportDates.Latest10k ||
		*company.Latest10q != reportDates.Latest10q {
		c.EnterCompanyInfo(company.Cik, email)
		c.EnterFinancialReports(secRes, email)

		if companyErr != nil {
			log.Fatal(companyErr)
		}

		if company == nil {
			log.Fatalf("No company data found in DB for CIK: %d", cik)
		}

	}
}

func (c *Controller) GetCikByTickerC(ticker string) string {
	cik, cikErr := c.repo.GetCikByTickerR(ticker)
	if cikErr != nil {
		log.Fatal(cikErr)
	}

	return cik
}
