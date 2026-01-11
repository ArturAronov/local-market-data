package company_info

import (
	"encoding/json"
	"log"
	"market-data/src/utils"
)

func GetCompanyTickersC(email string) {
	url := "https://www.sec.gov/files/company_tickers.json"

	body, bodyErr := utils.HttpReq(email, url)
	if bodyErr != nil {
		log.Fatalf("[GetCompanyTickersC] failed to handle request %v", bodyErr)
	}

	var secRes SecResponse
	jsonErr := json.Unmarshal(body, &secRes)
	if jsonErr != nil {
		log.Fatalf("[GetCompanyTickersC] Failed to unmarshal JSON:\n%v\n", jsonErr)
	}

	InsertTickerInfoR(&secRes)

	// return res.StatusCode, &secRes, nil
}
