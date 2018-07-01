package conversation

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/fmitra/dennis-bot/internal/actions"
	"github.com/fmitra/dennis-bot/pkg/alphapoint"
	"github.com/fmitra/dennis-bot/pkg/crypto"
	"github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/users"
	mocks "github.com/fmitra/dennis-bot/test"
)

type OnboardUserSuite struct {
	suite.Suite
	Env    *mocks.TestEnv
	Action *actions.Actions
}

func (suite *OnboardUserSuite) SetupSuite() {
	configFile := "../../config/config.json"
	suite.Env = mocks.GetTestEnv(configFile)
	suite.Action = &actions.Actions{
		Db:         suite.Env.Db,
		Cache:      suite.Env.Cache,
		Config:     suite.Env.Config,
		Alphapoint: &alphapoint.Client{},
	}
}

func (suite *OnboardUserSuite) TearDownSuite() {
	mocks.CleanUpEnv(suite.Env)
}

func (suite *OnboardUserSuite) BeforeTest(suiteName, testName string) {
	MessageMap = mocks.MessageMapMock
	mocks.CleanUpEnv(suite.Env)
}

func (suite *OnboardUserSuite) TestGetResponseList() {
	onboardUser := &OnboardUser{}
	assert.Equal(suite.T(), 6, len(onboardUser.GetResponses()))
}

func (suite *OnboardUserSuite) TestAsksForPassword() {
	onboardUser := &OnboardUser{
		&Conversation{
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
		&Conversation{
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
		&Conversation{
			Step:       2,
			IncMessage: incMessage,
		},
		suite.Action,
	}

	response, err := onboardUser.ValidatePassword()
	assert.Equal(suite.T(), BotResponse("Okay try again later"), response)
	assert.EqualError(suite.T(), err, "password rejected")

	message = mocks.GetMockMessage("YES")
	json.Unmarshal(message, &incMessage)

	onboardUser = &OnboardUser{
		&Conversation{
			Step:       2,
			IncMessage: incMessage,
		},
		suite.Action,
	}

	response, err = onboardUser.ValidatePassword()
	assert.Equal(suite.T(), BotResponse("I'm having trouble with this password"), response)
	assert.EqualError(suite.T(), err, "cipher text too short")

	message = mocks.GetMockMessage("Invalid")
	json.Unmarshal(message, &incMessage)

	onboardUser = &OnboardUser{
		&Conversation{
			Step:       2,
			IncMessage: incMessage,
		},
		suite.Action,
	}

	response, err = onboardUser.ValidatePassword()
	assert.Equal(suite.T(), BotResponse("I didn't understand that"), response)
	assert.EqualError(suite.T(), err, "response invalid")

	message = mocks.GetMockMessage("YES")
	json.Unmarshal(message, &incMessage)
	password, _ := crypto.Encrypt("password", suite.Env.Config.SecretKey)

	onboardUser = &OnboardUser{
		&Conversation{
			Step:       2,
			IncMessage: incMessage,
			AuxData:    password,
		},
		suite.Action,
	}

	response, err = onboardUser.ValidatePassword()
	assert.Equal(suite.T(), BotResponse(""), response)
	assert.NoError(suite.T(), err)
}

func (suite *OnboardUserSuite) TestCreatesUser() {
	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("Yes")
	json.Unmarshal(message, &incMessage)
	password, _ := crypto.Encrypt("password", suite.Env.Config.SecretKey)

	onboardUser := &OnboardUser{
		&Conversation{
			Step:       2,
			IncMessage: incMessage,
			AuxData:    password,
		},
		suite.Action,
	}

	response, err := onboardUser.ValidatePassword()

	var user users.User
	suite.Env.Db.Where("telegram_id = ?", mocks.TestUserID).First(&user)

	assert.Equal(suite.T(), user.TelegramID, mocks.TestUserID)
	assert.Equal(suite.T(), BotResponse(""), response)
	assert.NoError(suite.T(), err)
}

func (suite *OnboardUserSuite) TestReturnsErrorForFailedAccountCreation() {
	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("Yes")
	json.Unmarshal(message, &incMessage)
	mocks.CreateTestUser(suite.Env.Db, 0)
	password, _ := crypto.Encrypt("password", suite.Env.Config.SecretKey)

	onboardUser := &OnboardUser{
		&Conversation{
			Step:       2,
			IncMessage: incMessage,
			AuxData:    password,
		},
		suite.Action,
	}

	response, err := onboardUser.ValidatePassword()

	assert.EqualError(suite.T(), err, "account creation failed")
	assert.Equal(suite.T(), BotResponse("Couldn't create account"), response)
}

func (suite *OnboardUserSuite) TestAsksForCurrency() {
	onboardUser := &OnboardUser{
		&Conversation{},
		suite.Action,
	}
	response, err := onboardUser.AskForCurrency()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), BotResponse("What currency do you want to use?"), response)
}

func (suite *OnboardUserSuite) TestValidatesCurrency() {
	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("abc")
	json.Unmarshal(message, &incMessage)
	mocks.CreateTestUser(suite.Env.Db, 0)

	onboardUser := &OnboardUser{
		&Conversation{
			IncMessage: incMessage,
		},
		suite.Action,
	}

	response, err := onboardUser.ValidateCurrency()
	assert.EqualError(suite.T(), err, "invalid currency")
	assert.Equal(suite.T(), BotResponse("Currency is invalid"), response)

	message = mocks.GetMockMessage("SGD")
	json.Unmarshal(message, &incMessage)

	onboardUser = &OnboardUser{
		&Conversation{
			IncMessage: incMessage,
		},
		suite.Action,
	}

	response, err = onboardUser.ValidateCurrency()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), BotResponse(""), response)
}

func (suite *OnboardUserSuite) TestSaysOutro() {
	onboardUser := &OnboardUser{
		&Conversation{
			Step: 5,
		},
		suite.Action,
	}

	response, _ := onboardUser.SayOutro()
	assert.Equal(suite.T(), BotResponse("Outro message"), response)
}

func TestOnboardUserSuite(t *testing.T) {
	suite.Run(t, new(OnboardUserSuite))
}
