package company_facts

import (
	"fmt"
	"log"
	"market-data/src/user"
	"market-data/src/utils"
)

func GetCompanyFactsC(cik int) {
	cikStr := fmt.Sprintf("%010d", cik)
	url := fmt.Sprintf("https://data.sec.gov/submissions/CIK%s.json", cikStr)

	email, emailErr := user.GetUserEmail()
	if emailErr != nil {
		log.Fatalf("[GetCompanyFactsC] Error getting user email: %v\n", emailErr)
	}

	body, _, bodyErr := utils.HttpReq(*email, url)
	if bodyErr != nil {
		log.Fatalf("[GetCompanyTickersC] failed to handle request %v", bodyErr)
	}

	fmt.Println(string(body))
}
