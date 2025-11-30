package company_ticker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type secEntry struct {
	CIKStr int64  `json:"cik_str"`
	Ticker string `json:"ticker"`
	Title  string `json:"title"`
}

type secResponse map[string]secEntry

func GetCompanyTickersC(email string) (int, *secResponse, error) {
	client := &http.Client{}
	url := "https://www.sec.gov/files/company_tickers.json"

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		log.Fatalf("[GetCompanyTickers] Could not create request: %v\n", reqErr)
	}

	req.Header.Set("User-Agent", email)

	res, respErr := client.Do(req)
	if respErr != nil {
		log.Fatalf("[GetCompanyTickers] Error making http request: %v\n", respErr)
	}

	defer res.Body.Close()

	body, bodyErr := io.ReadAll(res.Body)
	if bodyErr != nil {
		log.Fatalf("[GetCompanyTickers] Could not read response body: %v\n", bodyErr)
	}

	if res.StatusCode > 499 {
		err := fmt.Errorf(
			"[GetCompanyTickers] Failed to fetch GetCompanyTickers\nStatus code: %d\nbody:body: %v\n",
			res.StatusCode,
			string(body),
		)

		return res.StatusCode, nil, err
	}

	var secRes secResponse
	if res.StatusCode < 300 {
		jsonErr := json.Unmarshal(body, &secRes)
		if jsonErr != nil {
			log.Fatalf("[GetCompanyTickers] Failed to unmarshal JSON:\n%v\n", jsonErr)
		}
	}

	InsertTickerInfoR(&secRes)

	return res.StatusCode, &secRes, nil
}
