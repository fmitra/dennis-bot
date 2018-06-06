package internal

import (
	"fmt"
	"strconv"

	"github.com/fmitra/dennis-bot/pkg/alphapoint"
	"github.com/fmitra/dennis-bot/pkg/expenses"
	"github.com/fmitra/dennis-bot/pkg/wit"
)

type Actions struct {
	env         *Env
	witResponse wit.WitResponse
	userId      int
}

// An intent represents a user's objective. We perform a separate
// action for each supported intent and retreive an appropriate response
// for the action performed.
func (a *Actions) ProcessIntent() BotResponse {
	intent := a.witResponse.GetIntent()
	var botResponse BotResponse

	switch intent {
	case wit.INTENT_TRACKING_SUCCESS:
		botResponse = a.TrackExpense()
	case wit.INTENT_PERIOD_TOTAL_SUCCESS:
		botResponse = a.GetExpenseTotal()
	default:
		botResponse = a.GetDefaultMessage()
	}

	return botResponse
}

// Logs a new expense entry for the user
func (a *Actions) createNewExpense() bool {
	date := a.witResponse.GetDate()
	amount, fromCurrency, _ := a.witResponse.GetAmount()
	targetCurrency := "USD"
	description, _ := a.witResponse.GetDescription()

	var conversion alphapoint.Conversion
	var newConversion *alphapoint.Conversion
	cacheKey := fmt.Sprintf("%s_%s", fromCurrency, targetCurrency)
	a.env.cache.Get(cacheKey, &conversion)

	historicalAmount := conversion.Rate * amount
	if conversion.Rate == 0 {
		ap := alphapoint.NewClient(a.env.config.AlphaPoint.Token)
		historicalAmount, newConversion = ap.Convert(
			fromCurrency,
			"USD",
			amount,
		)
		a.env.cache.Set(cacheKey, newConversion)
	}

	expense := &expenses.Expense{
		Date:        date,
		Description: description,
		Total:       amount,
		Historical:  historicalAmount,
		Currency:    fromCurrency,
		UserId:      a.userId,
	}
	expenseManager := expenses.NewExpenseManager(a.env.db)
	return expenseManager.Save(expense)
}

// Starts a goroutine to track an expense and returns an appropriate message
func (a *Actions) TrackExpense() BotResponse {
	var messageVar string

	go a.createNewExpense()
	return GetMessage(TRACKING_SUCCESS, messageVar)
}

// Retrieves historical expense total by time period and returns a message
func (a *Actions) GetExpenseTotal() BotResponse {
	expenseManager := expenses.NewExpenseManager(a.env.db)
	period, err := a.witResponse.GetSpendPeriod()
	total, err := expenseManager.TotalByPeriod(period, a.userId)
	messageVar := strconv.FormatFloat(total, 'f', 2, 64)

	if err != nil {
		return GetMessage(ERROR, messageVar)
	}

	return GetMessage(PERIOD_TOTAL_SUCCESS, messageVar)
}

// Returns a default message for an unsupported user intent
func (a *Actions) GetDefaultMessage() BotResponse {
	var messageVar string
	return GetMessage(DEFAULT, messageVar)
}
