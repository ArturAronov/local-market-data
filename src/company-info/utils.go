package company_info

import (
	company_metadata "market-data/src/company-metadata"
	"strings"
)

func EnterCompanyInfo(cik int) error {
	var companyData Company
	submissions := company_metadata.GetSubmissionsC(cik)
	reportDates := company_metadata.GetLatestReportDate(submissions.Filings.Recent)

	companyData.Cik = cik
	companyData.Sic = submissions.SIKStr
	companyData.Name = submissions.Name
	companyData.Phone = submissions.Phone
	companyData.EntryType = submissions.EntryType
	companyData.OwnerOrg = submissions.OwnerOrg
	companyData.Exchanges = strings.Join(submissions.Exchanges, ",")
	companyData.Description = submissions.Description
	companyData.FiscalYearEnd = submissions.FiscalYearEnd

	if reportDates != nil {
		companyData.Latest10k = reportDates.Latest10k
		companyData.Latest10q = reportDates.Latest10q
	}

	updateErr := UpdateCompanyR(companyData)
	if updateErr != nil {
		return updateErr
	}

	return nil
}
