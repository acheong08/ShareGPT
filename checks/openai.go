// This file contains checks for the validity of OpenAI API keys and how much credit is left.
package checks

import (
	"encoding/json"
	"net/http"

	"github.com/acheong08/ShareGPT/typings"
)

func GetCredits(apiKey string) (typings.BillingSubscription, error) {
	// Make request
	req, err := http.NewRequest("GET", "https://api.openai.com/dashboard/billing/subscription", nil)
	if err != nil {
		return typings.BillingSubscription{}, err
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")
	// Send request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return typings.BillingSubscription{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return typings.BillingSubscription{}, err
	}
	// Parse response
	var creditSummary typings.BillingSubscription
	// Map response to struct
	err = json.NewDecoder(response.Body).Decode(&creditSummary)
	if err != nil {
		return typings.BillingSubscription{}, err
	}
	return creditSummary, nil
}

func GetGrants(apiKey string) (typings.CreditSummary, error) {
	req, err := http.NewRequest("GET", "https://api.openai.com/dashboard/billing/credit_grants", nil)
	if err != nil {
		return typings.CreditSummary{}, err
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")
	// Send request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return typings.CreditSummary{}, err
	}
	defer response.Body.Close()
	var creditSummary typings.CreditSummary
	// Map response to struct
	err = json.NewDecoder(response.Body).Decode(&creditSummary)
	if err != nil {
		return typings.CreditSummary{}, err
	}
	if response.StatusCode != 200 {
		return creditSummary, err
	}
	return creditSummary, nil
}

func GetTotalCredits(apiKey string) (float64, error) {
	creditSummary, err := GetGrants(apiKey)
	if err != nil {
		return 0, err
	}
	var totalCredits float64 = 0
	totalCredits += creditSummary.TotalAvailable
	billingSummary, err := GetCredits(apiKey)
	if err != nil {
		return 0, err
	}
	totalCredits += billingSummary.HardLimitUSD
	return totalCredits, nil
}
