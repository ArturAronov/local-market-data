package utils

type ComapnyInfoT struct {
	Cik        int                          `json:"cik"`
	EntityName string                       `json:"entityName"`
	Facts      map[FactKeysT]map[string]any `json:"facts"`
}

func ExtractFactKeysToMap(body ComapnyInfoT, factKey FactKeysT, target map[string]bool) {
	facts := body.Facts[factKey]
	for k := range facts {
		target[k] = true
	}
}
