package scripts

import (
	"encoding/json"
	"fmt"
	"log"
	"market-data/utils"
	"os"
	"path/filepath"
	"strings"
)

func stringSetToSlice(set map[string]bool) []byte {
	result := make([]string, 0, len(set))
	for key := range set {
		result = append(result, key)
	}

	return []byte(strings.Join(result, "\n"))
}

func FileReader() {
	files, _ := os.ReadDir("test-data/")

	factsMapDei := make(map[string]bool)
	factsMapSrt := make(map[string]bool)
	factsMapIfrs := make(map[string]bool)
	factsMapInvest := make(map[string]bool)
	factsMapUsGaap := make(map[string]bool)

	for i, file := range files {
		rawFile := fmt.Sprintf("%v%v", "test-data/", file)
		parsedFile := strings.ReplaceAll(rawFile, "- ", "")

		fileBody, fileBodyErr := os.ReadFile(parsedFile)
		if fileBodyErr != nil {
			log.Fatal(fileBodyErr)
		}
		var bodyJson utils.ComapnyInfoT
		json.Unmarshal(fileBody, &bodyJson)

		pct := float32(i) / float32(len(files)) * 100
		fmt.Printf("%.2f%% >> Getting data for %s \n", pct, bodyJson.EntityName)

		utils.ExtractFactKeysToMap(bodyJson, utils.Dei, factsMapDei)
		utils.ExtractFactKeysToMap(bodyJson, utils.Srt, factsMapSrt)
		utils.ExtractFactKeysToMap(bodyJson, utils.Ifrs, factsMapIfrs)
		utils.ExtractFactKeysToMap(bodyJson, utils.Invest, factsMapInvest)
		utils.ExtractFactKeysToMap(bodyJson, utils.UsGaap, factsMapUsGaap)
	}

	deiOutPath := filepath.Join("facts/", fmt.Sprintf("%s.md", utils.Dei))
	srtOutPath := filepath.Join("facts/", fmt.Sprintf("%s.md", utils.Srt))
	ifrsOutPath := filepath.Join("facts/", fmt.Sprintf("%s.md", utils.Ifrs))
	investOutPath := filepath.Join("facts/", fmt.Sprintf("%s.md", utils.Invest))
	usGaapOutPath := filepath.Join("facts/", fmt.Sprintf("%s.md", utils.UsGaap))

	if writeErr := os.WriteFile(deiOutPath, stringSetToSlice(factsMapSrt), 0o644); writeErr != nil {
		// log.Printf("[GetCompanyFacts] write %s: %v. \n", outPath, writeErr)
		return
	}
	if writeErr := os.WriteFile(srtOutPath, stringSetToSlice(factsMapSrt), 0o644); writeErr != nil {
		// log.Printf("[GetCompanyFacts] write %s: %v. \n", outPath, writeErr)
		return
	}
	if writeErr := os.WriteFile(ifrsOutPath, stringSetToSlice(factsMapIfrs), 0o644); writeErr != nil {
		// log.Printf("[GetCompanyFacts] write %s: %v. \n", outPath, writeErr)
		return
	}
	if writeErr := os.WriteFile(investOutPath, stringSetToSlice(factsMapInvest), 0o644); writeErr != nil {
		// log.Printf("[GetCompanyFacts] write %s: %v. \n", outPath, writeErr)
		return
	}
	if writeErr := os.WriteFile(usGaapOutPath, stringSetToSlice(factsMapUsGaap), 0o644); writeErr != nil {
		// log.Printf("[GetCompanyFacts] write %s: %v. \n", outPath, writeErr)
		return
	}
}
