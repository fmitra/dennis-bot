package conversation

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/wit"
	mocks "github.com/fmitra/dennis-bot/test"
)

type TrackExpenseSuite struct {
	suite.Suite
	Env    *mocks.TestEnv
	Action *Actions
}

func (suite *TrackExpenseSuite) SetupSuite() {
	configFile := "../../config/config.json"
	suite.Env = mocks.GetTestEnv(configFile)
	suite.Action = &Actions{
		suite.Env.Db,
		suite.Env.Cache,
		suite.Env.Config,
	}
}

func (suite *TrackExpenseSuite) BeforeTest(suiteName, testName string) {
	MessageMap = mocks.MessageMapMock
}

func (suite *TrackExpenseSuite) AfterTest(suiteName, testName string) {
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
	var witResponse wit.WitResponse
	json.Unmarshal(rawWitResponse, &witResponse)

	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("")
	json.Unmarshal(message, &incMessage)

	trackExpense := &TrackExpense{
		Context{
			Step:        0,
			WitResponse: witResponse,
			IncMessage:  incMessage,
		},
		suite.Action,
	}
	response, cx := trackExpense.Respond()
	assert.Equal(suite.T(), BotResponse("Roger that!"), response)
	assert.Equal(suite.T(), -1, cx.Step)
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
	var witResponse wit.WitResponse
	json.Unmarshal(rawWitResponse, &witResponse)

	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("")
	json.Unmarshal(message, &incMessage)

	trackExpense := &TrackExpense{
		Context{
			Step:        0,
			WitResponse: witResponse,
			IncMessage:  incMessage,
		},
		suite.Action,
	}
	response, cx := trackExpense.Respond()
	assert.Equal(suite.T(), BotResponse("Whoops!"), response)
	assert.Equal(suite.T(), -1, cx.Step)
}

func TestTrackExpenseSuite(t *testing.T) {
	suite.Run(t, new(TrackExpenseSuite))
}
