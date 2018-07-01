package conversation

import (
	"errors"
	"fmt"
	"strconv"

	a "github.com/fmitra/dennis-bot/internal/actions"
	"github.com/fmitra/dennis-bot/pkg/crypto"
	"github.com/fmitra/dennis-bot/pkg/sessions"
	"github.com/fmitra/dennis-bot/pkg/users"
)

// GetExpenseTotal is an Intent designed to retrieve expense history totals
type GetExpenseTotal struct {
	*Conversation
	actions *a.Actions
}

// GetResponses returns a list of functions each containing a BotResponse.
func (i *GetExpenseTotal) GetResponses() []func() (BotResponse, error) {
	return []func() (BotResponse, error){
		i.AskForPassword,
		i.ValidatePassword,
		i.CalculateTotal,
	}
}

// AskForPassword requests a user for their password.
func (i *GetExpenseTotal) AskForPassword() (BotResponse, error) {
	expensePeriod, err := i.WitResponse.GetSpendPeriod()
	if err != nil {
		return GetMessage(GetExpenseTotalInvalidPeriod, ""), errors.New("invalid period")
	}

	i.AuxData = expensePeriod
	telegramUserID := i.IncMessage.GetUser().ID

	key := i.actions.Config.SecretKey
	_, err = passwordInCache(telegramUserID, key, i.actions.Cache)
	if err == nil {
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
	key := i.actions.Config.SecretKey

	_, err := passwordInCache(telegramUserID, key, i.actions.Cache)
	if err == nil {
		return i.SkipResponse()
	}

	shouldEnd := i.IncMessage.GetMessage() == "cancel"
	if shouldEnd {
		i.EndConversation()
		return GetMessage(GetExpenseTotalCancel, ""), errors.New("user requested cancel")
	}

	password := i.IncMessage.GetMessage()
	encryptedPass, err := crypto.Encrypt(password, i.actions.Config.SecretKey)
	if err != nil {
		response := GetMessage(GetExpenseTotalPasswordInvalid, "")
		return response, err
	}

	manager := users.NewUserManager(i.actions.Db)
	user := manager.GetByTelegramID(telegramUserID)

	if err = user.ValidatePassword(password); err != nil {
		response := GetMessage(GetExpenseTotalPasswordInvalid, "")
		err = errors.New("password invalid")
		return response, err
	}

	threeMinutes := 180
	cacheKey := passwordCacheKey(telegramUserID)
	i.actions.Cache.Set(cacheKey, encryptedPass, threeMinutes)
	return i.SkipResponse()
}

// CalculateTotal runs an action to check for the total sum of user expenses
// for a specific period in time.
func (i *GetExpenseTotal) CalculateTotal() (BotResponse, error) {
	telegramUserID := i.IncMessage.GetUser().ID
	key := i.actions.Config.SecretKey

	// No need to handle this error or below. Bot will return an error
	// response if the private key is invalid
	password, _ := passwordInCache(telegramUserID, key, i.actions.Cache)

	manager := users.NewUserManager(i.actions.Db)
	user := manager.GetByTelegramID(telegramUserID)
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
func passwordInCache(userID uint, secretKey string, cache sessions.Session) (string, error) {
	var password string
	cacheKey := passwordCacheKey(userID)
	err := cache.Get(cacheKey, &password)
	if err != nil {
		return "", err
	}

	decryptedPass, err := crypto.Decrypt(password, secretKey)
	if err != nil {
		return "", err
	}

	return decryptedPass, nil
}
