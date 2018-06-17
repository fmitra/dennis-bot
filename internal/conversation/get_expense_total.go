package conversation

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/fmitra/dennis-bot/pkg/users"
)

// A Conversation designed to retrieve expense history totals
type GetExpenseTotal struct {
	Context
	actions *Actions
}

func (i *GetExpenseTotal) GetResponses() []func() (BotResponse, error) {
	return []func() (BotResponse, error){
		i.AskForPassword,
		i.ValidatePassword,
		i.CalculateTotal,
	}
}

func (i *GetExpenseTotal) Respond() (BotResponse, *Context) {
	responses := i.GetResponses()
	return i.Process(responses)
}

func (i *GetExpenseTotal) AskForPassword() (BotResponse, error) {
	expensePeriod, err := i.WitResponse.GetSpendPeriod()
	if err != nil {
		return GetMessage(GET_EXPENSE_TOTAL_INVALID_PERIOD, ""), errors.New("Invalid period")
	}

	i.AuxData = expensePeriod
	telegramUserId := i.IncMessage.GetUser().Id
	cacheKey := fmt.Sprintf("%s_password", strconv.Itoa(int(telegramUserId)))
	var password string
	err = i.actions.Cache.Get(cacheKey, &password)
	if err == nil {
		return i.SkipResponse()
	}

	return GetMessage(GET_EXPENSE_TOTAL_ASK_FOR_PASSWORD, ""), nil
}

func (i *GetExpenseTotal) ValidatePassword() (BotResponse, error) {
	telegramUserId := i.IncMessage.GetUser().Id
	cacheKey := fmt.Sprintf("%s_password", strconv.Itoa(int(telegramUserId)))
	var password string
	err := i.actions.Cache.Get(cacheKey, &password)
	if err == nil {
		return i.SkipResponse()
	}

	shouldStop := i.IncMessage.GetMessage() == "stop"
	if shouldStop {
		i.EndConversation()
		return BotResponse(""), errors.New("User requested to stop")
	}

	password = i.IncMessage.GetMessage()
	manager := users.NewUserManager(i.actions.Db)
	user := manager.GetByTelegramId(telegramUserId)
	if !user.IsPasswordValid(password) {
		response := GetMessage(GET_EXPENSE_TOTAL_PASSWORD_INVALID, "")
		err := errors.New("Password invalid")
		return response, err
	}

	threeMinutes := 180
	i.actions.Cache.Set(cacheKey, password, threeMinutes)
	return i.SkipResponse()
}

func (i *GetExpenseTotal) CalculateTotal() (BotResponse, error) {
	var password string
	telegramUserId := i.IncMessage.GetUser().Id
	cacheKey := fmt.Sprintf("%s_password", strconv.Itoa(int(telegramUserId)))
	i.actions.Cache.Get(cacheKey, &password)

	manager := users.NewUserManager(i.actions.Db)
	user := manager.GetByTelegramId(telegramUserId)
	// No need to handle this error. Bot will return an error response
	// if the private key is invalid
	privateKey, _ := user.GetPrivateKey(password)

	expensePeriod := i.AuxData
	messageVar, err := i.actions.GetExpenseTotal(expensePeriod, i.BotUserId, privateKey)
	var response BotResponse

	response = GetMessage(GET_EXPENSE_TOTAL_SUCCESS, messageVar)
	if err != nil {
		response = GetMessage(GET_EXPENSE_TOTAL_ERROR, "")
	}

	i.EndConversation()
	return response, nil
}
