package company_facts

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	company_ticker "market-data/src/company-tickers"
	"market-data/src/user"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type SecEntry struct {
	CIKStr int64  `json:"cik_str"`
	Ticker string `json:"ticker"`
	Title  string `json:"title"`
}

const SLEEP_BETWEEN = 1 * time.Second

func GetCompanyFacts() {
	email, err := user.GetUserEmail()
	if err != nil {
		log.Fatalf("[GetCompanyFacts] get user email: %v", err)
	}

	entries, err := company_ticker.GetTickerInfoR()
	if err != nil {
		log.Fatalf("[GetCompanyFacts] get tickers: %v", err)
	}

	outDir := "test-data"
	if mkErr := os.MkdirAll(outDir, 0o755); mkErr != nil {
		log.Fatalf("[GetCompanyFacts] mkdir %s: %v", outDir, mkErr)
	}

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	for _, entry := range entries {
		if entry.CIKStr == 0 {
			continue
		}

		cik10 := fmt.Sprintf("%010d", entry.CIKStr)
		url := fmt.Sprintf("https://data.sec.gov/api/xbrl/companyfacts/CIK%s.json", cik10)

		req, reqErr := http.NewRequest("GET", url, nil)
		if reqErr != nil {
			log.Printf("[GetCompanyFacts] new request CIK %s: %v", cik10, reqErr)
			continue
		}

		// SEC guidance: app name/version + contact email
		req.Header.Set("User-Agent", fmt.Sprintf("tes-fetch/1.0 (%s)", *email))
		req.Header.Set("Accept", "application/json")

		res, doErr := client.Do(req)
		if doErr != nil {
			log.Printf("[GetCompanyFacts] http CIK %s: %v", cik10, doErr)
			// time.Sleep(SLEEP_BETWEEN)
			continue
		}

		func() {
			defer res.Body.Close()

			if res.StatusCode == http.StatusNotFound {
				log.Printf("[GetCompanyFacts] CIK %v not found, deleting it from company_tickers. \n", entry.CIKStr)
				company_ticker.DeleteTickerR(entry.CIKStr)

				return
			}

			if res.StatusCode != http.StatusOK {
				snippet, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
				log.Printf("[GetCompanyFacts] non-200 CIK %s: %s body=%s, \n",
					cik10, res.Status, string(snippet))
				return
			}

			body, readErr := io.ReadAll(res.Body)
			if readErr != nil {
				log.Printf("[GetCompanyFacts] read body CIK %s: %v. \n", cik10, readErr)
				return
			}

			// Optional: validate JSON so we don’t persist HTML/error
			var js any
			if err := json.Unmarshal(body, &js); err != nil {
				log.Printf("[GetCompanyFacts] invalid JSON CIK %s: %v. \n", cik10, err)
				return
			}

			outPath := filepath.Join(outDir, fmt.Sprintf("%s.json", cik10))
			if writeErr := os.WriteFile(outPath, body, 0o644); writeErr != nil {
				log.Printf("[GetCompanyFacts] write %s: %v. \n", outPath, writeErr)
				return
			}

			log.Printf("[GetCompanyFacts] saved %s (%s %s). \n", outPath, entry.Ticker, entry.Title)
		}()

		time.Sleep(SLEEP_BETWEEN)
	}
}
