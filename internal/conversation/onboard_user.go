package conversation

import (
	"errors"
	"strings"
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

func (i *OnboardUser) Respond() (BotResponse, int) {
	responses := i.GetResponses()
	return i.Process(responses)
}

func (i *OnboardUser) AskForPassword() (BotResponse, error) {
	messageVar := ""
	return GetMessage(ONBOARD_USER_ASK_FOR_PASSWORD, messageVar), nil
}

func (i *OnboardUser) ConfirmPassword() (BotResponse, error) {
	messageVar := i.IncMessage.GetMessage()
	return GetMessage(ONBOARD_USER_CONFIRM_PASSWORD, messageVar), nil
}

func (i *OnboardUser) ValidatePassword() (BotResponse, error) {
	messageVar := ""
	userInput := strings.ToLower(i.IncMessage.GetMessage())
	var response BotResponse
	var err error

	if userInput == "no" {
		response = GetMessage(ONBOARD_USER_REJECT_PASSWORD, messageVar)
		err = errors.New("Password rejected")
		i.EndConversation()
	} else if userInput != "yes" {
		response = GetMessage(ONBOARD_USER_CONFIRM_PASSWORD_ERROR, messageVar)
		err = errors.New("Response invalid")
	}

	return response, err
}

func (c *OnboardUser) SayOutro() (BotResponse, error) {
	messageVar := ""
	return GetMessage(ONBOARD_USER_SAY_OUTRO, messageVar), nil
}
