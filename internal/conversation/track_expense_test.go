package conversation

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/fmitra/dennis-bot/internal/actions"
	"github.com/fmitra/dennis-bot/pkg/alphapoint"
	"github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/wit"
	mocks "github.com/fmitra/dennis-bot/test"
)

type TrackExpenseSuite struct {
	suite.Suite
	Env    *mocks.TestEnv
	Action *actions.Actions
}

func (suite *TrackExpenseSuite) SetupSuite() {
	configFile := "../../config/config.json"
	suite.Env = mocks.GetTestEnv(configFile)
	suite.Action = &actions.Actions{
		Db:         suite.Env.Db,
		Cache:      suite.Env.Cache,
		Config:     suite.Env.Config,
		Alphapoint: &alphapoint.Client{},
	}
}

func (suite *TrackExpenseSuite) TearDownSuite() {
	mocks.CleanUpEnv(suite.Env)
}

func (suite *TrackExpenseSuite) BeforeTest(suiteName, testName string) {
	MessageMap = mocks.MessageMapMock
	mocks.CleanUpEnv(suite.Env)
}

func (suite *TrackExpenseSuite) TestGetResponseList() {
	trackExpense := &TrackExpense{}
	assert.Equal(suite.T(), 1, len(trackExpense.GetResponses()))
}

func (suite *TrackExpenseSuite) TestReturnsSuccessMessage() {
	rawWitResponse := []byte(`{
		"entities": {
			"amount": [
				{ "value": "20 SGD", "confidence": 100.00 }
			],
			"datetime": [
				{ "value": "", "confidence": 100.00 }
			],
			"description": [
				{ "value": "Food", "confidence": 100.00 }
			]
		}
	}`)
	var witResponse wit.Response
	json.Unmarshal(rawWitResponse, &witResponse)

	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("")
	json.Unmarshal(message, &incMessage)

	trackExpense := &TrackExpense{
		&Conversation{
			Step:        0,
			WitResponse: witResponse,
			IncMessage:  incMessage,
		},
		suite.Action,
	}
	response, err := trackExpense.ConfirmExpense()
	assert.Equal(suite.T(), BotResponse("Roger that!"), response)
	assert.NoError(suite.T(), err)
}

func (suite *TrackExpenseSuite) TestReturnsErrorMessage() {
	rawWitResponse := []byte(`{
		"entities": {
			"amount": [],
			"datetime": [
				{ "value": "", "confidence": 100.00 }
			],
			"description": []
		}
	}`)
	var witResponse wit.Response
	json.Unmarshal(rawWitResponse, &witResponse)

	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("")
	json.Unmarshal(message, &incMessage)

	trackExpense := &TrackExpense{
		&Conversation{
			Step:        0,
			WitResponse: witResponse,
			IncMessage:  incMessage,
		},
		suite.Action,
	}
	response, err := trackExpense.ConfirmExpense()
	assert.Equal(suite.T(), BotResponse("Whoops!"), response)
	assert.NoError(suite.T(), err)
}

func TestTrackExpenseSuite(t *testing.T) {
	suite.Run(t, new(TrackExpenseSuite))
}
