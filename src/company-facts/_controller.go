// package company_facts

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"market-data/src/user"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

type CompanyData struct {
	Cik   int64                 `json:"cik"`
	Name  string                `json:"entityName"`
	Facts map[string]FactsGroup `json:"facts"`
}

type FactsGroup map[string]Label

type Label struct {
	Label       string          `json:"label"`
	Description string          `json:"description"`
	Units       json.RawMessage `json:"units"`
}

type Data struct {
	Start *string `json:"start"`
	End   string  `json:"end"`
	Val   float64 `json:"val"`
	Accn  string  `json:"accn"`
	Fy    int     `json:"fy"`
	Fp    string  `json:"fp"`
	Form  string  `json:"form"`
	Filed string  `json:"filed"`
	Frame *string `json:"frame"`
}

type ProcessedLabel struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	Data        []Data `json:"data"`
}

type ProcessedFacts map[string]map[string]ProcessedLabel

type ProcessedCompanyData struct {
	Cik   int64          `json:"cik"`
	Name  string         `json:"entityName"`
	Facts ProcessedFacts `json:"facts"`
}

func flattenUnits(raw json.RawMessage) ([]Data, error) {
	var units map[string][]Data
	unitsErr := json.Unmarshal(raw, &units)
	if unitsErr != nil {
		return nil, fmt.Errorf("failed to unmarshal units: %w", unitsErr)
	}

	var result []Data
	for _, dataList := range units {
		result = append(result, dataList...)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].End < result[j].End
	})

	return result, nil
}

func GetCompanyFacts(cik int64) *ProcessedCompanyData {
	client := &http.Client{}

	url := fmt.Sprintf("https://data.sec.gov/api/xbrl/companyfacts/CIK%010d.json", cik)

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		log.Fatalf("[GetCompanyFacts] Could not create request: %v\n", reqErr)
	}

	email, emailErr := user.GetUserEmail()
	if emailErr != nil {
		log.Fatalf("[GetCompanyFacts] Failed to get user email: %v\n", emailErr)
	}

	req.Header.Set("User-Agent", *email)

	res, respErr := client.Do(req)
	if respErr != nil {
		log.Fatalf("[GetCompanyFacts] Error making HTTP request: %v\n", respErr)
	}
	defer res.Body.Close()

	body, bodyErr := io.ReadAll(res.Body)
	if bodyErr != nil {
		log.Fatalf("[GetCompanyFacts] Could not read response body: %v\n", bodyErr)
	}

	var rawData CompanyData
	if jsonErr := json.Unmarshal(body, &rawData); jsonErr != nil {
		log.Fatalf("[GetCompanyFacts] Could not parse JSON: %v\n", jsonErr)
	}

	processedData := ProcessedCompanyData{
		Cik:   rawData.Cik,
		Name:  rawData.Name,
		Facts: make(ProcessedFacts),
	}

	for taxonomyName, factMap := range rawData.Facts {
		processedData.Facts[taxonomyName] = make(map[string]ProcessedLabel)

		for key, label := range factMap {
			dataPoints, err := flattenUnits(label.Units)
			if err != nil {
				log.Printf("Error flattening %s::%s: %v\n", taxonomyName, key, err)
				continue
			}

			processedData.Facts[taxonomyName][key] = ProcessedLabel{
				Label:       label.Label,
				Description: label.Description,
				Data:        dataPoints,
			}
		}
	}

	// Output to JSON file
	outputPath := filepath.Join(".", "test_all_taxonomies.json")
	outputFile, createErr := os.Create(outputPath)
	if createErr != nil {
		log.Fatalf("[GetCompanyFacts] Could not create output file: %v\n", createErr)
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ")

	if encodeErr := encoder.Encode(processedData); encodeErr != nil {
		log.Fatalf("[GetCompanyFacts] Could not write JSON to file: %v\n", encodeErr)
	}

	return &processedData
}
