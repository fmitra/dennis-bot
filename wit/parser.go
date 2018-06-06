package wit

import (
	"errors"
	"log"
	"time"

	"github.com/fmitra/dennis/utils"
)

const (
	INTENT_TRACKING_SUCCESS = "tracking_success"

	INTENT_TRACKING_ERROR = "tracking_error"

	INTENT_PERIOD_TOTAL_SUCCESS = "period_total_success"

	INTENT_UNKNOWN = "default"
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
		TotalSpent  WitEntity `json:"total_spent"`
	} `json:"entities"`
}

// Checks if Wit.ai was able to infer a total spent query
func (w WitResponse) GetSpendPeriod() (string, error) {
	totalSpent := w.Entities.TotalSpent
	if len(totalSpent) == 0 {
		return "", errors.New("No period specified")
	}

	return totalSpent[0].Value, nil
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
		log.Printf("wit: cannot infer without amount")
		return false, err
	}

	description, err := w.GetDescription()
	if err != nil {
		log.Printf("wit: cannot infer without description")
		return true, err
	}

	log.Printf("wit: inferring tracking %s %s %s", amount, currency, description)
	return true, nil
}

func (w WitResponse) IsRequestingTotal() (bool, error) {
	spendPeriod, err := w.GetSpendPeriod()
	if err != nil {
		log.Printf("wit: cannot infer spend period")
		return false, err
	}

	log.Printf("wit: inferring spend period %s", spendPeriod)
	return true, nil
}

func (w WitResponse) GetIntent() string {
	isTracking, trackingErr := w.IsTracking()
	isRequestingTotal, totalErr := w.IsRequestingTotal()

	if isTracking && trackingErr != nil {
		return INTENT_TRACKING_ERROR
	} else if isTracking && trackingErr == nil {
		return INTENT_TRACKING_SUCCESS
	} else if isRequestingTotal && totalErr == nil {
		return INTENT_PERIOD_TOTAL_SUCCESS
	} else {
		return INTENT_UNKNOWN
	}
}
