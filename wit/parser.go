package wit

import (
	"errors"
	"log"
	"time"

	"github.com/fmitra/dennis/utils"
)

// Wit.ai Entity
type WitEntity []struct {
	Value      string  `json:"value"`
	Confidence float64 `json:"confidence"`
}

// Wit.ai API Response
type WitResponse struct {
	Entities struct {
		Amount      WitEntity `json:"amount"`
		DateTime    WitEntity `json:"datetime"`
		Description WitEntity `json:"description"`
	} `json:"entities"`
}

// Checks if Wit.ai was able to infer a valid
// Amount Entity from the IncomingMessage.getMessage()
func (w WitResponse) GetAmount() (float64, string, error) {
	amount := w.Entities.Amount
	if len(amount) == 0 {
		return 0, "", errors.New("No amount")
	}

	totalAmount, currency := utils.ParseAmount(amount[0].Value)

	if totalAmount > 0 && currency != "" {
		return totalAmount, currency, nil
	}

	return 0, "", errors.New("Invalid amount")
}

// Checks if Wit.ai was able to infer a valid
// Description Entity from the IncomingMessage.getMessage()
func (w WitResponse) GetDescription() (string, error) {
	description := w.Entities.Description
	if len(description) == 0 {
		return "", errors.New("No description")
	}

	parsedDescription := utils.ParseDescription(description[0].Value)
	return parsedDescription, nil
}

// Checks if Wit.ai was able to infer a valid
// Date Entity from the IncomingMessage.getMessage()
// If no date is provided, it will always default to today
func (w WitResponse) GetDate() time.Time {
	dateTime := w.Entities.DateTime
	stringDate := ""
	if len(dateTime) != 0 {
		stringDate = dateTime[0].Value
	}

	parsedDate := utils.ParseDate(stringDate)
	return parsedDate
}

// Infers whether the User is attempting to track an expense
func (w WitResponse) IsTracking() (bool, error) {
	amount, currency, err := w.GetAmount()
	if err != nil {
		log.Printf("Cannot infer without amount")
		return false, err
	}

	description, err := w.GetDescription()
	if err != nil {
		log.Printf("Cannot infer without description")
		return true, err
	}

	log.Printf("Inferring tracking %s %s %s", amount, currency, description)
	return true, nil
}
