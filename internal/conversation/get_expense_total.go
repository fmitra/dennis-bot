package conversation

// A Conversation designed to retrieve expense history totals
type GetExpenseTotal struct {
	Context
	actions *Actions
}

func (i *GetExpenseTotal) GetResponses() []func() (BotResponse, error) {
	return []func() (BotResponse, error){
		i.CalculateTotal,
	}
}

func (i *GetExpenseTotal) Respond() (BotResponse, int) {
	responses := i.GetResponses()
	return i.Process(responses)
}

func (i *GetExpenseTotal) CalculateTotal() (BotResponse, error) {
	userId := i.IncMessage.GetUser().Id
	messageVar, err := i.actions.GetExpenseTotal(i.WitResponse, userId)
	var response BotResponse

	response = GetMessage(GET_EXPENSE_TOTAL_SUCCESS, messageVar)
	if err != nil {
		response = GetMessage(GET_EXPENSE_TOTAL_ERROR, "")
	}

	i.EndConversation()
	return response, nil
}