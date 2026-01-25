package company_info

import "time"

type SecEntry struct {
	CIKStr int64  `json:"cik_str"`
	Ticker string `json:"ticker"`
	Title  string `json:"title"`
}

type SecEntryRes map[string]SecEntry

// SEC company-facts api types
type ReportData struct {
	Start string  `json:"start"` // "2009-06-01",
	End   string  `json:"end"`   // "2010-08-31",
	Val   float64 `json:"val"`   // 39662816,
	Accn  string  `json:"accn"`  // "0001104659-10-049632",
	Fy    int     `json:"fy"`    // 2011,
	Fp    string  `json:"fp"`    // "Q1",
	Form  string  `json:"form"`  // "10-Q",
	Filed string  `json:"filed"` // "2010-09-23",
	Frame string  `json:"frame"` // "CY2010Q3I"
}

type Fact struct {
	Label       string                  `json:"label"`       // Adjustments to Additional Paid in Capital, Income Tax Benefit from Share-based Compensation
	Description string                  `json:"description"` // Amount of increase in additional paid in capital (APIC) resulting from a tax benefit associated with share-based compensation plan other than an employee stock ownership plan (ESOP). Includes, but is not limited to, excess tax benefit
	Units       map[string][]ReportData `json:"units"`       // USD
}

type CompanyFacts struct {
	Cik        int                        `json:"cik"`        // 88948
	EntityName string                     `json:"entityName"` // Seneca Foods Corporation
	Facts      map[string]map[string]Fact `json:"facts"`
}

// Company db table types
type DbCompany struct {
	Cik           int     `json:"cik"`
	Sic           string  `json:"sic"`
	Name          string  `json:"name"`
	Ticker        string  `json:"ticker"`
	Phone         string  `json:"phone"`
	EntryType     string  `json:"entry_type"`
	OwnerOrg      string  `json:"owner_org"`
	Exchanges     string  `json:"exchanges"`
	Description   string  `json:"description"`
	FiscalYearEnd string  `json:"fiscal_year_end"`
	Latest10k     *string `json:"latest_10k"`
	Latest10q     *string `json:"latest_10q"`
}

type DbFact struct {
	Cik         int
	FactKey     string
	Namespace   string
	Label       string
	Description string
	Unit        string
}

type DbReport struct {
	Id      int
	Cik     int
	FactKey string
	Start   *time.Time
	End     *time.Time
	Val     float64
	Accn    *string
	Fy      *int
	Fp      *string
	Form    string
	Filed   time.Time
	Frame   *string
	Hash    []byte
}
