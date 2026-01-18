package company_metadata

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
	CIKStr               string             `json:"cik"`                  // "0001018724"
	SIKStr               string             `json:"sik"`                  // "5961"
	Name                 string             `json:"name"`                 // "AMAZON COM INC",
	EntryType            string             `json:"entityType"`           // "operating"
	Description          string             `json:"sicDescription"`       // "sicDescription": "Retail-Catalog & Mail-Order Houses"
	OwnerOrg             string             `json:"ownerOrg"`             // "07 Trade & Services"
	Exchanges            []string           `json:"exchanges"`            // ["Nasdaq"]
	FiscalYearEnd        string             `json:"fiscalYearEnd"`        // "1231"
	StateOfIncorporation string             `json:"stateOfIncorporation"` // "DE"
	Phone                string             `json:"phone"`                // "2062661000"
	Addresses            map[string]Address `json:"addresses"`
	Filings              FilingsWrapper     `json:"filings"`
}
