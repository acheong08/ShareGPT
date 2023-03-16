package typings

type CreditSummary struct {
	Object         string      `json:"object"`
	TotalGranted   float64     `json:"total_granted"`
	TotalUsed      float64     `json:"total_used"`
	TotalAvailable float64     `json:"total_available"`
	Grants         CreditGrant `json:"grants"`
	Error          OpenAIError `json:"error"`
}

type CreditGrant struct {
	Object string        `json:"object"`
	Data   []CreditGrant `json:"data"`
}

type CreditGrantItem struct {
	Object      string  `json:"object"`
	ID          string  `json:"id"`
	GrantAmount float64 `json:"grant_amount"`
	UsedAmount  float64 `json:"used_amount"`
	EffectiveAt float64 `json:"effective_at"`
	ExpiresAt   float64 `json:"expires_at"`
}

type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	// Param can be null or a string
	Param string `json:"param"`
	Code  string `json:"code"`
}

type APIKeySubmission struct {
	APIKey string `json:"api_key"`
}
