package main

import (
	"context"
	"encoding/csv"
	"log"

	// company_facts "market-data/src/company-facts"
	"market-data/utils"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func csvTest(key, label, taxonomy, company string, cik int64) {
	fileName := "data.csv"

	_, statErr := os.Stat(fileName)
	fileExists := !os.IsNotExist(statErr)

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Printf("error closing file: %v", cerr)
		}
	}()

	w := csv.NewWriter(f)

	if !fileExists {
		if err := w.Write([]string{"key", "label", "taxonomy", "cik", "company"}); err != nil {
			log.Fatalf("error writing header to csv: %v", err)
		}
	}

	record := []string{key, label, taxonomy, strconv.FormatInt(cik, 10), company}
	if err := w.Write(record); err != nil {
		log.Fatalf("error writing record to csv: %v", err)
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatalf("csv writer error: %v", err)
	}
}

func main() {
	dbs, dbsErr := utils.InitDB()
	if dbsErr != nil {
		log.Fatalln(dbsErr)
	}

	userDb := dbs[utils.USER_EMAIL]
	marketDb := dbs[utils.COMPANY_TICKERS]

	defer userDb.Close()
	defer marketDb.Close()

	// data := company_facts.GetCompanyFacts(320193)

	// taxonomies := make([]string, 0, len(data.Facts))
	// for tax := range data.Facts {
	// 	taxonomies = append(taxonomies, tax)
	// }
	// sort.Strings(taxonomies)

	// for _, tax := range taxonomies {
	// 	keys := make([]string, 0, len(data.Facts[tax]))
	// 	for k := range data.Facts[tax] {
	// 		keys = append(keys, k)
	// 	}
	// 	sort.Strings(keys)

	// 	for _, k := range keys {
	// 		csvTest(k, data.Facts[tax][k].Label, tax, data.Name, data.Cik)
	// 	}
	// }

	// utils.ParamsResover()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

}
