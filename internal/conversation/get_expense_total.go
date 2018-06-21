package conversation

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/fmitra/dennis-bot/pkg/sessions"
	"github.com/fmitra/dennis-bot/pkg/users"
)

// GetExpenseTotal is an Intent designed to retrieve expense history totals
type GetExpenseTotal struct {
	Context
	actions *Actions
}

// GetResponses returns a list of functions each containing a BotResponse.
func (i *GetExpenseTotal) GetResponses() []func() (BotResponse, error) {
	return []func() (BotResponse, error){
		i.AskForPassword,
		i.ValidatePassword,
		i.CalculateTotal,
	}
}

// Respond proccesses a list of response functions.
func (i *GetExpenseTotal) Respond() (BotResponse, *Context) {
	responses := i.GetResponses()
	return i.Process(responses)
}

// AskForPassword requests a user for their password.
func (i *GetExpenseTotal) AskForPassword() (BotResponse, error) {
	expensePeriod, err := i.WitResponse.GetSpendPeriod()
	if err != nil {
		return GetMessage(GetExpenseTotalInvalidPeriod, ""), errors.New("invalid period")
	}

	i.AuxData = expensePeriod
	telegramUserID := i.IncMessage.GetUser().ID

	if err := passwordInCache(telegramUserID, i.actions.Cache); err == nil {
		return i.SkipResponse()
	}

	return GetMessage(GetExpenseTotalAskForPassword, ""), nil
}

// ValidatePassword checks if the supplied password in a previous message is correct.
// This is a validation response, therefore it will return an empty response
// nil error on success, triggering the bot to skip over to the next response function
// in line.
func (i *GetExpenseTotal) ValidatePassword() (BotResponse, error) {
	telegramUserID := i.IncMessage.GetUser().ID

	if err := passwordInCache(telegramUserID, i.actions.Cache); err == nil {
		return i.SkipResponse()
	}

	shouldStop := i.IncMessage.GetMessage() == "stop"
	if shouldStop {
		i.EndConversation()
		return BotResponse(""), errors.New("user requested to stop")
	}

	password := i.IncMessage.GetMessage()
	manager := users.NewUserManager(i.actions.Db)
	user := manager.GetByTelegramID(telegramUserID)

	if err := user.ValidatePassword(password); err != nil {
		response := GetMessage(GetExpenseTotalPasswordInvalid, "")
		err := errors.New("password invalid")
		return response, err
	}

	threeMinutes := 180
	cacheKey := passwordCacheKey(telegramUserID)
	i.actions.Cache.Set(cacheKey, password, threeMinutes)
	return i.SkipResponse()
}

// CalculateTotal runs an action to check for the total sum of user expenses
// for a specific period in time.
func (i *GetExpenseTotal) CalculateTotal() (BotResponse, error) {
	var password string
	telegramUserID := i.IncMessage.GetUser().ID
	cacheKey := passwordCacheKey(telegramUserID)
	i.actions.Cache.Get(cacheKey, &password)

	manager := users.NewUserManager(i.actions.Db)
	user := manager.GetByTelegramID(telegramUserID)
	// No need to handle this error. Bot will return an error response
	// if the private key is invalid
	privateKey, _ := user.GetPrivateKey(password)

	expensePeriod := i.AuxData
	messageVar, err := i.actions.GetExpenseTotal(expensePeriod, i.BotUserID, privateKey)
	if err != nil {
		return GetMessage(GetExpenseTotalError, ""), nil
	}

	i.EndConversation()
	return GetMessage(GetExpenseTotalSuccess, messageVar), nil
}

// passwordCacheKey returns the cache key for password checks.
func passwordCacheKey(userID uint) string {
	return fmt.Sprintf("%s_password", strconv.Itoa(int(userID)))
}

// passwordInCache checks for a password in cache, indicating the user completed
// this flow within the past few minutes. We do not require users to re-enter password
// for consecutive queries within the cache timeout.
func passwordInCache(userID uint, cache sessions.Session) error {
	var password string
	cacheKey := passwordCacheKey(userID)
	return cache.Get(cacheKey, &password)
}
