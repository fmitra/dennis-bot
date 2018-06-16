package conversation

import (
	"errors"
	"strings"

	"github.com/fmitra/dennis-bot/pkg/crypto"
)

// An Intent designed to onboard a new user into the bot platform
type OnboardUser struct {
	Context
	actions *Actions
}

func (i *OnboardUser) GetResponses() []func() (BotResponse, error) {
	return []func() (BotResponse, error){
		i.AskForPassword,
		i.ConfirmPassword,
		i.ValidatePassword,
		i.SayOutro,
	}
}

func (i *OnboardUser) Respond() (BotResponse, *Context) {
	responses := i.GetResponses()
	return i.Process(responses)
}

func (i *OnboardUser) AskForPassword() (BotResponse, error) {
	messageVar := ""
	return GetMessage(ONBOARD_USER_ASK_FOR_PASSWORD, messageVar), nil
}

func (i *OnboardUser) ConfirmPassword() (BotResponse, error) {
	password := i.IncMessage.GetMessage()
	hashedPassword, err := crypto.HashText(password)
	if err != nil {
		return GetMessage(ONBOARD_USER_PASSWORD_HASH_FAILED, ""), err
	}

	i.AuxData = hashedPassword
	return GetMessage(ONBOARD_USER_CONFIRM_PASSWORD, password), nil
}

func (i *OnboardUser) ValidatePassword() (BotResponse, error) {
	messageVar := ""
	userInput := strings.ToLower(i.IncMessage.GetMessage())
	var isUserCreated bool
	var response BotResponse
	var err error

	isPasswordConfirmed := userInput == "yes"
	isPasswordRejected := userInput == "no"

	if isPasswordRejected {
		response = GetMessage(ONBOARD_USER_REJECT_PASSWORD, messageVar)
		err = errors.New("Password rejected")
		i.EndConversation()
		return response, err
	}

	if isPasswordConfirmed {
		// Password should have been set to auxiliary data in the previous step
		password := i.AuxData
		userId := i.IncMessage.GetUser().Id
		isUserCreated = i.actions.CreateNewUser(userId, password)
	}

	if isPasswordConfirmed && isUserCreated {
		return BotResponse(""), nil
	}

	if isPasswordConfirmed && !isUserCreated {
		response = GetMessage(ONBOARD_USER_ACCOUNT_CREATION_FAILED, messageVar)
		err = errors.New("Account creation failed")
		return response, err
	}

	response = GetMessage(ONBOARD_USER_CONFIRM_PASSWORD_ERROR, messageVar)
	err = errors.New("Response invalid")
	return response, err
}

func (i *OnboardUser) SayOutro() (BotResponse, error) {
	i.EndConversation()
	messageVar := ""
	return GetMessage(ONBOARD_USER_SAY_OUTRO, messageVar), nil
}
