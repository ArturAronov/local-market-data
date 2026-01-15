package submission_data

type RecentFilings struct {
	FilingDate []string `json:"filingDate"`
	Form       []string `json:"form"`
}

type FilingsWrapper struct {
	Recent RecentFilings `json:"recent"`
}

type Address struct {
	Street1                   string `json:"street1"`
	Street2                   string `json:"street2"`
	City                      string `json:"city"`
	StateOrCountry            string `json:"stateOrCountry"`
	ZipCode                   string `json:"zipCode"`
	StateOrCountryDescription string `json:"stateOrCountryDescription"`
	IsForeignLocation         int    `json:"isForeignLocation"`
	ForeignStateTerritory     string `json:"foreignStateTerritory"`
	Country                   string `json:"country"`
	CountryCode               string `json:"countryCode"`
}

type SubmissionData struct {
	CIKStr               string             `json:"cik"`
	SIKStr               string             `json:"sik"`
	EntryType            string             `json:"entityType"`
	Description          string             `json:"sicDescription"`
	OwnerOrg             string             `json:"ownerOrg"`
	Exchanges            []string           `json:"exchanges"`
	FiscalYearEnd        string             `json:"fiscalYearEnd"`
	StateOfIncorporation string             `json:"stateOfIncorporation"`
	Addresses            map[string]Address `json:"addresses"`
	Phone                string             `json:"phone"`
	Filings              FilingsWrapper     `json:"filings"`
}
