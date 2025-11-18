package company_ticker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type companyInfo struct {
	Ticker string `json:"ticker"`
	Title  string `json:"title"`
}

type secEntry struct {
	CIKStr int64  `json:"cik_str"`
	Ticker string `json:"ticker"`
	Title  string `json:"title"`
}

type CIKMap map[string]companyInfo
type secResponse map[string]secEntry

func GetCompanyTickers() (int, *CIKMap, error) {
	client := &http.Client{}
	url := "https://www.sec.gov/files/company_tickers.json"

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		log.Fatalf("GetCompanyTickers: Could not create request: %v\n", reqErr)
	}

	req.Header.Set("User-Agent", "adsf@asd.ee")

	resp, respErr := client.Do(req)
	if respErr != nil {
		log.Fatalf("GetCompanyTickers: Error making http request: %v\n", respErr)
	}

	defer resp.Body.Close()

	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		log.Fatalf("GetCompanyTickers: Could not read response body: %v\n", bodyErr)
	}

	if resp.StatusCode > 399 {
		err := fmt.Errorf(
			"GetCompanyTickers: Failed to fetch GetCompanyTickers\nStatus code: %d\nbody:body: %v\n",
			resp.StatusCode,
			string(body),
		)

		return resp.StatusCode, nil, err
	}

	var secResp secResponse
	jsonErr := json.Unmarshal(body, &secResp)
	if jsonErr != nil {
		log.Fatalf("GetCompanyTickers: Failed to unmarshal JSON:\n%v\n", jsonErr)
	}

	cikMap := make(CIKMap)
	for _, entry := range secResp {
		cik := fmt.Sprintf("%010d", entry.CIKStr)
		cikMap[cik] = companyInfo{
			Ticker: entry.Ticker,
			Title:  entry.Title,
		}
	}

	return resp.StatusCode, &cikMap, nil
}
