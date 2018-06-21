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

type GenericResponseSuite struct {
	suite.Suite
	Env *mocks.TestEnv
}

func (suite *GenericResponseSuite) BeforeTest(suiteName, testName string) {
	MessageMap = mocks.MessageMapMock
}

func (suite *GenericResponseSuite) TestGetResponseList() {
	genericResponse := &GenericResponse{}
	assert.Equal(suite.T(), 1, len(genericResponse.GetResponses()))
}

func (suite *GenericResponseSuite) TestReturnResponse() {
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

	genericResponse := &GenericResponse{
		Context{
			Step:        0,
			WitResponse: witResponse,
			IncMessage:  incMessage,
		},
		&Actions{},
	}

	response, cx := genericResponse.Respond()
	assert.Equal(suite.T(), -1, cx.Step)
	assert.Equal(suite.T(), BotResponse("This is a default message"), response)
}

func TestGenericResponseSuite(t *testing.T) {
	suite.Run(t, new(GenericResponseSuite))
}
