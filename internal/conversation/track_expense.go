package conversation

import (
	"github.com/fmitra/dennis-bot/pkg/wit"

	"github.com/fmitra/dennis-bot/pkg/users"
)

// TrackExpense is an Intent designed to track a user's expenses.
type TrackExpense struct {
	Context
	actions *Actions
}

// GetResponses proccesses a list of response functions.
func (i *TrackExpense) GetResponses() []func() (BotResponse, error) {
	return []func() (BotResponse, error){
		i.ConfirmExpense,
	}
}

// Respond proccesses a list of response functions.
func (i *TrackExpense) Respond() (BotResponse, *Context) {
	responses := i.GetResponses()
	return i.Process(responses)
}

// ConfirmExpense starts an action to track the user's expense and returns
// confirmation if it was successful or if it failed.
func (i *TrackExpense) ConfirmExpense() (BotResponse, error) {
	var messageVar string
	var response BotResponse
	overview := i.WitResponse.GetMessageOverview()

	telegramUserID := i.IncMessage.GetUser().ID
	manager := users.NewUserManager(i.actions.Db)
	user := manager.GetByTelegramID(telegramUserID)
	publicKey, _ := user.GetPublicKey()

	response = GetMessage(TrackExpenseError, messageVar)
	if overview == wit.TrackingRequestedSuccess {
		go i.actions.CreateNewExpense(i.WitResponse, i.BotUserID, publicKey)
		response = GetMessage(TrackExpenseSuccess, messageVar)
	}

	i.EndConversation()
	return response, nil
}
