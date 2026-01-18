package company_metadata

import (
	"encoding/json"
	"fmt"
	"log"
	"market-data/src/user"
	"market-data/src/utils"
)

func GetSubmissionsC(cik int) *SubmissionData {
	cikStr := fmt.Sprintf("%010d", cik)
	url := fmt.Sprintf("https://data.sec.gov/submissions/CIK%s.json", cikStr)

	email, emailErr := user.GetUserEmail()
	if emailErr != nil {
		log.Fatalf("[GetSubmissionsC] Error getting user email: %v\n", emailErr)
	}

	body, _, bodyErr := utils.HttpReq(*email, url)
	if bodyErr != nil {
		log.Fatalf("[GetSubmissionsC] failed to handle request %v", bodyErr)
	}

	var secRes SubmissionData
	if err := json.Unmarshal(body, &secRes); err != nil {
		log.Fatalf("[GetSubmissionsC] unmarshal failed: %v\n", err)
	}

	return &secRes
}
