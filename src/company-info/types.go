package company_info

type SecEntry struct {
	CIKStr int64  `json:"cik_str"`
	Ticker string `json:"ticker"`
	Title  string `json:"title"`
}

type SecEntryRes map[string]SecEntry
