package wit

import (
	"errors"
	"log"
	"time"

	"github.com/fmitra/dennis-bot/pkg/utils"
)

const (
	// TrackingRequestedSuccess indicates the user is requesting to track an expense
	TrackingRequestedSuccess = "tracking_requested_success"

	// TrackingRequestedError indicates the user is requested to track an expense
	// but failed to provide the necessary data behind the expense
	TrackingRequestedError = "tracking_requested_error"

	// ExpenseTotalRequestedSuccess indicates the user is attempting to
	// get a sum of their expense history
	ExpenseTotalRequestedSuccess = "expense_total_requested_success"

	// UnknownRequest indicates Wit.ai failed to infer context around a message
	UnknownRequest = "unknown_request"
)

// Entity is an item Wit.ai inferred from a response.
// All Entities have a Confidence property indicating an
// estimate of Wit.ai's inference.
type Entity []struct {
	Value      string  `json:"value"`
	Confidence float64 `json:"confidence"`
}

// Response is a a Response from Wit.ai containing a payload of Entities.
type Response struct {
	Entities struct {
		Amount      Entity `json:"amount"`
		DateTime    Entity `json:"datetime"`
		Description Entity `json:"description"`
		TotalSpent  Entity `json:"total_spent"`
	} `json:"entities"`
}

// GetSpendPeriod returns the spending period a user requested, for example,
// a user may be interested in one `month` of spending history.
func (r *Response) GetSpendPeriod() (string, error) {
	totalSpent := r.Entities.TotalSpent
	if len(totalSpent) == 0 {
		return "", errors.New("no period specified")
	}

	return totalSpent[0].Value, nil
}

// GetAmount returns the total amount a user is trying to track.
func (r *Response) GetAmount() (float64, string, error) {
	amount := r.Entities.Amount
	if len(amount) == 0 {
		return 0, "", errors.New("no amount")
	}

	totalAmount, currency := utils.ParseAmount(amount[0].Value)

	if totalAmount > 0 && currency != "" {
		return totalAmount, currency, nil
	}

	return 0, "", errors.New("invalid amount")
}

// GetDescription returns the description of an expense.
func (r *Response) GetDescription() (string, error) {
	description := r.Entities.Description
	if len(description) == 0 {
		return "", errors.New("no description")
	}

	parsedDescription := utils.ParseDescription(description[0].Value)
	return parsedDescription, nil
}

// GetDate returns the date of an expense. If no date is provided,
// we default to today.
func (r *Response) GetDate() time.Time {
	dateTime := r.Entities.DateTime
	stringDate := ""
	if len(dateTime) != 0 {
		stringDate = dateTime[0].Value
	}

	parsedDate := utils.ParseDate(stringDate)
	return parsedDate
}

// IsTracking infers whether the user is trying to track an expense.
func (r Response) IsTracking() (bool, error) {
	_, _, err := r.GetAmount()
	if err != nil {
		return false, err
	}

	_, err = r.GetDescription()
	if err != nil {
		log.Printf("wit: cannot infer tracking without description")
		return true, err
	}

	return true, nil
}

// IsRequestingTotal infers whether the user is requesting a sum
// of their expense history.
func (r Response) IsRequestingTotal() (bool, error) {
	_, err := r.GetSpendPeriod()
	if err != nil {
		log.Printf("wit: cannot infer expense total without spend period")
		return false, err
	}

	return true, nil
}

// GetMessageOverview returns a a description of what Wit.ai
// inferred from a user's message.
func (r Response) GetMessageOverview() string {
	isTracking, trackingErr := r.IsTracking()
	isRequestingTotal, totalErr := r.IsRequestingTotal()

	if isTracking && trackingErr != nil {
		return TrackingRequestedError
	} else if isTracking && trackingErr == nil {
		return TrackingRequestedSuccess
	} else if isRequestingTotal && totalErr == nil {
		return ExpenseTotalRequestedSuccess
	}

	return UnknownRequest
}
