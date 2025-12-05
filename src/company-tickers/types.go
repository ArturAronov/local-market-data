package company_ticker

type SecEntry struct {
	CIKStr int64  `json:"cik_str"`
	Ticker string `json:"ticker"`
	Title  string `json:"title"`
}

type SecResponse map[string]SecEntry
