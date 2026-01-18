package company_metadata

type ReportDates struct {
	Latest10k string
	Latest10q string
}

func GetLatestReportDate(data RecentFilings) *ReportDates {
	var index10k *int
	var index10q *int

	for i := range data.Form {
		if data.Form[i] == "10-K" {
			index10k = &i
		}

		if data.Form[i] == "10-Q" {
			index10q = &i
		}

		if index10k != nil && index10q != nil {
			break
		}
	}

	if index10k == nil && index10q == nil {
		return nil
	} else {
		return &ReportDates{
			Latest10k: data.FilingDate[*index10k],
			Latest10q: data.FilingDate[*index10q],
		}
	}
}
