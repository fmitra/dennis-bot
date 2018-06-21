package conversation

import (
	"errors"
	"strings"
)

// OnboardUser is an Intent designed to onboard a new user into the bot platform.
type OnboardUser struct {
	Context
	actions *Actions
}

// GetResponses returns a list of functions each containing a BotResponse.
func (i *OnboardUser) GetResponses() []func() (BotResponse, error) {
	return []func() (BotResponse, error){
		i.AskForPassword,
		i.ConfirmPassword,
		i.ValidatePassword,
		i.SayOutro,
	}
}

// Respond proccesses a list of response functions.
func (i *OnboardUser) Respond() (BotResponse, *Context) {
	responses := i.GetResponses()
	return i.Process(responses)
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
	i.AuxData = password
	return GetMessage(OnboardUserConfirmPassword, password), nil
}

// ValidatePassword validates a user's response from ConfirmPassword. It will
// return an empty response and nil error on success, triggering the bot
// to skip over to the next response function in line.
func (i *OnboardUser) ValidatePassword() (BotResponse, error) {
	messageVar := ""
	userInput := strings.ToLower(i.IncMessage.GetMessage())
	var response BotResponse
	var userCreatedError error
	var err error

	isPasswordConfirmed := userInput == "yes"
	isPasswordRejected := userInput == "no"

	if isPasswordRejected {
		response = GetMessage(OnboardUserRejectPassword, messageVar)
		err = errors.New("password rejected")
		i.EndConversation()
		return response, err
	}

	if isPasswordConfirmed {
		// Password should have been set to auxiliary data in the previous step
		password := i.AuxData
		userID := i.IncMessage.GetUser().ID
		userCreatedError = i.actions.CreateNewUser(userID, password)
	}

	if isPasswordConfirmed && userCreatedError == nil {
		return i.SkipResponse()
	}

	if isPasswordConfirmed && userCreatedError != nil {
		response = GetMessage(OnboardUserAccountCreationFailed, messageVar)
		err = errors.New("account creation failed")
		return response, err
	}

	response = GetMessage(OnboardUserConfirmPasswordError, messageVar)
	err = errors.New("response invalid")
	return response, err
}

// SayOutro confirms to the user that account creation is complete.
func (i *OnboardUser) SayOutro() (BotResponse, error) {
	i.EndConversation()
	messageVar := ""
	return GetMessage(OnboardUserSayOutro, messageVar), nil
}
