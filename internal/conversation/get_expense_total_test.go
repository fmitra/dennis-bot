package conversation

import (
	"encoding/json"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/users"
	"github.com/fmitra/dennis-bot/pkg/wit"
	mocks "github.com/fmitra/dennis-bot/test"
)

type ExpenseTotalSuite struct {
	suite.Suite
	Env    *mocks.TestEnv
	Action *Actions
}

func (suite *ExpenseTotalSuite) SetupSuite() {
	configFile := "../../config/config.json"
	suite.Env = mocks.GetTestEnv(configFile)
	suite.Action = &Actions{
		suite.Env.Db,
		suite.Env.Cache,
		suite.Env.Config,
	}
}

func (suite *ExpenseTotalSuite) BeforeTest(suiteName, testName string) {
	MessageMap = mocks.MessageMapMock
}

func (suite *ExpenseTotalSuite) AfterTest(suiteName, testName string) {
	mocks.CleanUpEnv(suite.Env)
}

func GetTestUser(db *gorm.DB) users.User {
	var user users.User
	db.Where("telegram_id = ?", mocks.TestUserId).First(&user)
	return user
}

func (suite *ExpenseTotalSuite) TestGetResponseList() {
	expenseTotal := &GetExpenseTotal{}
	assert.Equal(suite.T(), 1, len(expenseTotal.GetResponses()))
}

func (suite *ExpenseTotalSuite) TestGetExpenseTotalMessage() {
	rawWitResponse := []byte(`{
		"entities": {
			"amount": [],
			"datetime": [],
			"description": [],
			"total_spent": [
				{ "value": "month", "confidence": 100.00 }
			]
		}
	}`)
	var witResponse wit.WitResponse
	json.Unmarshal(rawWitResponse, &witResponse)

	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("")
	json.Unmarshal(message, &incMessage)

	expenseTotal := &GetExpenseTotal{
		Context{
			Step:        0,
			WitResponse: witResponse,
			IncMessage:  incMessage,
		},
		suite.Action,
	}
	response, step := expenseTotal.Respond()
	assert.Equal(suite.T(), BotResponse("You spent 0.00"), response)
	assert.Equal(suite.T(), -1, step)
}

func (suite *ExpenseTotalSuite) TestGetExpenseTotalError() {
	rawWitResponse := []byte(`{
		"entities": {
			"amount": [],
			"datetime": [],
			"description": [],
			"total_spent": [
				{ "value": "foo", "confidence": 100.00 }
			]
		}
	}`)
	var witResponse wit.WitResponse
	json.Unmarshal(rawWitResponse, &witResponse)

	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("")
	json.Unmarshal(message, &incMessage)

	expenseTotal := &GetExpenseTotal{
		Context{
			Step:        0,
			WitResponse: witResponse,
			IncMessage:  incMessage,
		},
		suite.Action,
	}
	response, step := expenseTotal.Respond()
	assert.Equal(suite.T(), BotResponse("Whoops!"), response)
	assert.Equal(suite.T(), -1, step)
}

func TestExpenseTotalSuite(t *testing.T) {
	suite.Run(t, new(ExpenseTotalSuite))
}
