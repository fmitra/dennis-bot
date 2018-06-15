package conversation

import (
	"github.com/fmitra/dennis-bot/pkg/wit"
)

// An Intent designed to track a user's expenses
type TrackExpense struct {
	Context
	actions *Actions
}

func (i *TrackExpense) GetResponses() []func() (BotResponse, error) {
	return []func() (BotResponse, error){
		i.ConfirmExpense,
	}
}

func (i *TrackExpense) Respond() (BotResponse, int) {
	responses := i.GetResponses()
	return i.Process(responses)
}

// Sends a success or error response for an expense tracking request
// based on Wit.ai's parsing.
func (i *TrackExpense) ConfirmExpense() (BotResponse, error) {
	messageVar := ""
	overview := i.WitResponse.GetMessageOverview()
	var response BotResponse

	response = GetMessage(TRACK_EXPENSE_ERROR, messageVar)
	if overview == wit.TRACKING_REQUESTED_SUCCESS {
		go i.actions.CreateNewExpense(i.WitResponse, i.BotUserId)
		response = GetMessage(TRACK_EXPENSE_SUCCESS, messageVar)
	}

	i.EndConversation()
	return response, nil
}
