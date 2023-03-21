package typings

import "time"

type BillingSubscription struct {
	Object             string          `json:"object"`
	HasPaymentMethod   bool            `json:"has_payment_method"`
	Canceled           bool            `json:"canceled"`
	CanceledAt         *time.Time      `json:"canceled_at"`
	Delinquent         *bool           `json:"delinquent"`
	AccessUntil        int64           `json:"access_until"`
	SoftLimit          int             `json:"soft_limit"`
	HardLimit          int             `json:"hard_limit"`
	SystemHardLimit    int             `json:"system_hard_limit"`
	SoftLimitUSD       float64         `json:"soft_limit_usd"`
	HardLimitUSD       float64         `json:"hard_limit_usd"`
	SystemHardLimitUSD float64         `json:"system_hard_limit_usd"`
	Plan               Plan            `json:"plan"`
	AccountName        string          `json:"account_name"`
	PONumber           *string         `json:"po_number"`
	BillingEmail       *string         `json:"billing_email"`
	TaxIDs             *[]string       `json:"tax_ids"`
	BillingAddress     BillingAddress  `json:"billing_address"`
	BusinessAddress    BusinessAddress `json:"business_address"`
}

type Plan struct {
	Title string `json:"title"`
	ID    string `json:"id"`
}

type BusinessAddress struct {
	City       string `json:"city"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

type BillingAddress struct {
	City       string  `json:"city"`
	Line1      string  `json:"line1"`
	Line2      *string `json:"line2"`
	State      string  `json:"state"`
	Country    string  `json:"country"`
	PostalCode string  `json:"postal_code"`
}

type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type APIKeySubmission struct {
	APIKey string `json:"api_key"`
}

type CreditSummary struct {
	Object         string      `json:"object"`
	TotalGranted   float64     `json:"total_granted"`
	TotalUsed      float64     `json:"total_used"`
	TotalAvailable float64     `json:"total_available"`
	Grants         CreditGrant `json:"grants"`
	Error          OpenAIError `json:"error"`
}

type CreditGrant struct {
	Object string            `json:"object"`
	Data   []CreditGrantItem `json:"data"`
}

type CreditGrantItem struct {
	Object      string  `json:"object"`
	ID          string  `json:"id"`
	GrantAmount float64 `json:"grant_amount"`
	UsedAmount  float64 `json:"used_amount"`
	EffectiveAt float64 `json:"effective_at"`
	ExpiresAt   float64 `json:"expires_at"`
}
