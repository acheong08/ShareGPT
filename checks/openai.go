// This file contains checks for the validity of OpenAI API keys and how much credit is left.
package checks

import (
	"encoding/json"
	"net/http"

	"github.com/acheong08/ShareGPT/typings"
)

func GetCredits(apiKey string) (typings.CreditSummary, error) {
	// Make request
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
	// Parse response
	var creditSummary typings.CreditSummary
	// Map response to struct
	err = json.NewDecoder(response.Body).Decode(&creditSummary)
	if err != nil {
		return typings.CreditSummary{}, err
	}
	return creditSummary, nil
}
