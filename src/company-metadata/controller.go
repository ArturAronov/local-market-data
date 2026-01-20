package company_metadata

import (
	"encoding/json"
	"fmt"
	"log"
	"market-data/src/utils"
)

func GetSubmissionsC(cik int, email string) *SubmissionData {
	cikStr := fmt.Sprintf("%010d", cik)
	url := fmt.Sprintf("https://data.sec.gov/submissions/CIK%s.json", cikStr)

	body, _, bodyErr := utils.HttpReq(email, url)
	if bodyErr != nil {
		log.Fatalf("[GetSubmissionsC] failed to handle request %v", bodyErr)
	}

	var secRes SubmissionData
	if err := json.Unmarshal(body, &secRes); err != nil {
		log.Fatalf("[GetSubmissionsC] unmarshal failed: %v\n", err)
	}

	return &secRes
}
