package conversation

import (
	"errors"
	"strings"

	a "github.com/fmitra/dennis-bot/internal/actions"
	"github.com/fmitra/dennis-bot/pkg/crypto"
	"github.com/fmitra/dennis-bot/pkg/users"
)

// OnboardUser is an Intent designed to onboard a new user into the bot platform.
type OnboardUser struct {
	*Conversation
	actions *a.Actions
}

// GetResponses returns a list of functions each containing a BotResponse.
func (i *OnboardUser) GetResponses() []func() (BotResponse, error) {
	return []func() (BotResponse, error){
		i.AskForPassword,
		i.ConfirmPassword,
		i.ValidatePassword,
		i.AskForCurrency,
		i.ValidateCurrency,
		i.SayOutro,
	}
}

// AskForPassword requests the User to create a new password in order to
// to create their account.
func (i *OnboardUser) AskForPassword() (BotResponse, error) {
	messageVar := ""
	return GetMessage(OnboardUserAskForPassword, messageVar), nil
}

// ConfirmPassword requests the user to confirm if the password they previously
// submitted is final.
func (i *OnboardUser) ConfirmPassword() (BotResponse, error) {
	password := i.IncMessage.GetMessage()
	encryptedPass, err := crypto.Encrypt(password, i.actions.Config.SecretKey)
	if err != nil {
		return GetMessage(OnboardUserConfirmPasswordInvalid, ""), err
	}

	i.AuxData = encryptedPass
	return GetMessage(OnboardUserConfirmPassword, password), nil
}

// ValidatePassword validates a user's response from ConfirmPassword. It will
// return an empty response and nil error on success, triggering the bot
// to skip over to the next response function in line.
func (i *OnboardUser) ValidatePassword() (BotResponse, error) {
	messageVar := ""
	userInput := strings.ToLower(i.IncMessage.GetMessage())
	var response BotResponse
	var err error

	isPasswordConfirmed := userInput == "yes"
	isPasswordRejected := userInput == "no"
	isInvalidResponse := !isPasswordConfirmed && !isPasswordRejected

	if isInvalidResponse {
		response = GetMessage(OnboardUserConfirmPasswordError, messageVar)
		err = errors.New("response invalid")
		return response, err
	}

	if isPasswordRejected {
		response = GetMessage(OnboardUserRejectPassword, messageVar)
		err = errors.New("password rejected")
		i.EndConversation()
		return response, err
	}

	// Password should have been set to auxiliary data in the previous step
	password, err := crypto.Decrypt(i.AuxData, i.actions.Config.SecretKey)
	if err != nil {
		response = GetMessage(OnboardUserDecryptionFailed, messageVar)
		i.EndConversation()
		return response, err
	}

	userID := i.IncMessage.GetUser().ID
	err = i.actions.CreateNewUser(userID, password)
	if err != nil {
		response = GetMessage(OnboardUserAccountCreationFailed, messageVar)
		err = errors.New("account creation failed")
		return response, err
	}

	return i.SkipResponse()
}

// AskForCurrency checks if a user would like to receive expense totals in a
// a currency other than the default (USD).
func (i *OnboardUser) AskForCurrency() (BotResponse, error) {
	return GetMessage(OnboardUserAskForCurrency, ""), nil
}

// ValidateCurrency checks if a user supplied a valid currency and updates
// the user's related Settings.
func (i *OnboardUser) ValidateCurrency() (BotResponse, error) {
	currency := i.IncMessage.GetMessage()
	tID := i.IncMessage.GetUser().ID
	manager := users.NewUserManager(i.actions.Db)
	user := manager.GetByTelegramID(tID)

	err := i.actions.SetUserCurrency(user.ID, currency)
	if err != nil {
		return GetMessage(OnboardUserInvalidCurrency, ""), err
	}

	return i.SkipResponse()
}

// SayOutro confirms to the user that account creation is complete.
func (i *OnboardUser) SayOutro() (BotResponse, error) {
	i.EndConversation()
	messageVar := ""
	return GetMessage(OnboardUserSayOutro, messageVar), nil
}
