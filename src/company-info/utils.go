package company_info

import (
	"log"
	company_metadata "market-data/src/company-metadata"
	"market-data/src/utils"
	"strings"
)

func (c *Controller) EnterCompanyInfo(cik int, email string) error {
	var companyData DbCompany
	submissions := company_metadata.GetSubmissionsC(cik, email)
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

	updateErr := c.repo.UpdateCompanyR(companyData)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (c *Controller) EnterCompanyFacts(data CompanyFacts) error {
	var factsData []DbFact
	var reportData []DbReport

	for namespace, factsByName := range data.Facts {
		for factName, fact := range factsByName {
			for factUnitName, reports := range fact.Units {
				fact := DbFact{
					Cik:         data.Cik,
					FactKey:     factName,
					Namespace:   namespace,
					Label:       fact.Label,
					Description: fact.Label,
					Unit:        factUnitName,
				}
				factsData = append(factsData, fact)

				for _, value := range reports {
					report := DbReport{
						Cik:     data.Cik,
						FactKey: factName,
						Val:     value.Val,
						Accn:    &value.Accn,
						Fy:      &value.Fy,
						Fp:      &value.Fp,
						Form:    value.Form,
						Frame:   &value.Frame,
					}

					if value.Start != "" {
						startDateSlice := strings.Split(value.Start, "-")

						startDate, startDateErr := utils.DateParser(utils.Date{
							Year:  startDateSlice[0],
							Month: startDateSlice[1],
							Day:   startDateSlice[2],
						})

						if startDateErr != nil {
							log.Panicf("EnterCompanyFacts: %v", startDateErr)
						}

						report.Start = &startDate
					}

					if value.End != "" {
						endDateSlice := strings.Split(value.End, "-")

						endDate, endDateErr := utils.DateParser(utils.Date{
							Year:  endDateSlice[0],
							Month: endDateSlice[1],
							Day:   endDateSlice[2],
						})

						if endDateErr != nil {
							log.Panicf("EnterCompanyFacts: %v", endDateErr)
						}

						report.End = &endDate
					}

					if value.Filed != "" {
						filedDateSlice := strings.Split(value.Filed, "-")

						filedDate, filedDateErr := utils.DateParser(utils.Date{
							Year:  filedDateSlice[0],
							Month: filedDateSlice[1],
							Day:   filedDateSlice[2],
						})

						if filedDateErr != nil {
							log.Panicf("EnterCompanyFacts: %v", filedDateErr)
						}

						report.Filed = filedDate
					}

					reportData = append(reportData, report)
				}
			}
		}
	}

	insertFactsErr := c.repo.InsertCompanyFactsR(factsData)
	if insertFactsErr != nil {
		return insertFactsErr
	}

	insertReportsErr := c.repo.InsertCompanyReportsR(reportData)
	if insertReportsErr != nil {
		return insertReportsErr
	}

	return nil
}
