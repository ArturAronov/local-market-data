package company_info

import (
	"encoding/json"
	"fmt"
	"log"

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

func (c *Controller) GetCompanyFactsC(cik int, email string) {
	cikStr := fmt.Sprintf("%010d", cik)
	url := fmt.Sprintf("https://data.sec.gov/api/xbrl/companyfacts/CIK%s.json", cikStr)

	body, _, bodyErr := utils.HttpReq(email, url)
	if bodyErr != nil {
		log.Fatalf("[GetCompanyFactsC] failed to handle request %v", bodyErr)
	}

	var secRes CompanyFacts
	jsonErr := json.Unmarshal(body, &secRes)
	if jsonErr != nil {
		log.Fatalf("[GetCompanyFactsC] failed to unmarshal JSON:\n%v\n,", jsonErr)
	}

	factCount, factCountErr := c.repo.CountCompanyFactsR(cik)
	if factCountErr != nil {
		log.Fatal(factCountErr)
	}

	if factCount == 0 {
		c.EnterCompanyInfo(cik, email)
		c.EnterCompanyFacts(secRes)
	} else {
		company, companyErr := c.repo.GetCompanyR(cik)
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
