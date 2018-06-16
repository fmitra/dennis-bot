package conversation

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/fmitra/dennis-bot/pkg/telegram"
	mocks "github.com/fmitra/dennis-bot/test"
)

type OnboardUserSuite struct {
	suite.Suite
	Env    *mocks.TestEnv
	Action *Actions
}

func (suite *OnboardUserSuite) SetupSuite() {
	configFile := "../../config/config.json"
	suite.Env = mocks.GetTestEnv(configFile)
	suite.Action = &Actions{
		suite.Env.Db,
		suite.Env.Cache,
		suite.Env.Config,
	}
}

func (suite *OnboardUserSuite) BeforeTest(suiteName, testName string) {
	MessageMap = mocks.MessageMapMock
}

func (suite *OnboardUserSuite) AfterTest(suiteName, testName string) {
	mocks.CleanUpEnv(suite.Env)
}

func (suite *OnboardUserSuite) TestGetResponseList() {
	onboardUser := &OnboardUser{}
	assert.Equal(suite.T(), 4, len(onboardUser.GetResponses()))
}

func (suite *OnboardUserSuite) TestAsksForPassword() {
	onboardUser := &OnboardUser{
		Context{
			Step: 0,
		},
		suite.Action,
	}

	response, _ := onboardUser.AskForPassword()
	assert.Equal(suite.T(), BotResponse("What's your password?"), response)
}

func (suite *OnboardUserSuite) TestConfirmsPassword() {
	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("")
	json.Unmarshal(message, &incMessage)

	onboardUser := &OnboardUser{
		Context{
			Step:       1,
			IncMessage: incMessage,
		},
		suite.Action,
	}

	response, _ := onboardUser.ConfirmPassword()
	assert.Equal(suite.T(), BotResponse("Your password is Hello world"), response)
}

func (suite *OnboardUserSuite) TestValidatesPassword() {
	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("No")
	json.Unmarshal(message, &incMessage)

	onboardUser := &OnboardUser{
		Context{
			Step:       2,
			IncMessage: incMessage,
		},
		suite.Action,
	}

	response, err := onboardUser.ValidatePassword()
	assert.Equal(suite.T(), BotResponse("Okay try again later"), response)
	assert.EqualError(suite.T(), err, "Password rejected")

	message = mocks.GetMockMessage("YES")
	json.Unmarshal(message, &incMessage)

	onboardUser = &OnboardUser{
		Context{
			Step:       2,
			IncMessage: incMessage,
		},
		suite.Action,
	}

	response, err = onboardUser.ValidatePassword()
	assert.Equal(suite.T(), BotResponse(""), response)
	assert.NoError(suite.T(), err)

	message = mocks.GetMockMessage("Invalid")
	json.Unmarshal(message, &incMessage)

	onboardUser = &OnboardUser{
		Context{
			Step:       2,
			IncMessage: incMessage,
		},
		suite.Action,
	}

	response, err = onboardUser.ValidatePassword()
	assert.Equal(suite.T(), BotResponse("I didn't understand that"), response)
	assert.EqualError(suite.T(), err, "Response invalid")
}

func (suite *OnboardUserSuite) TestSaysOutro() {
	onboardUser := &OnboardUser{
		Context{
			Step: 3,
		},
		suite.Action,
	}

	response, _ := onboardUser.SayOutro()
	assert.Equal(suite.T(), BotResponse("Outro message"), response)
}

func TestOnboardUserSuite(t *testing.T) {
	suite.Run(t, new(OnboardUserSuite))
}
