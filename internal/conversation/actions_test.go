package conversation

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/fmitra/dennis-bot/pkg/alphapoint"
	"github.com/fmitra/dennis-bot/pkg/wit"
	mocks "github.com/fmitra/dennis-bot/test"
)

type ActionSuite struct {
	suite.Suite
	Env    *mocks.TestEnv
	Action *Actions
}

func (suite *ActionSuite) SetupSuite() {
	configFile := "../../config/config.json"
	suite.Env = mocks.GetTestEnv(configFile)
	suite.Action = &Actions{
		suite.Env.Db,
		suite.Env.Cache,
		suite.Env.Config,
	}
}

func (suite *ActionSuite) AfterTest(suiteName, testName string) {
	mocks.CleanUpEnv(suite.Env)
}

func (suite *ActionSuite) TestCreatesNewExpense() {
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
	alphapointResponse := `{
		"Realtime Currency Exchange Rate": {
			"5. Exchange Rate": ".7"
		}
	}`
	alphapointServer := mocks.MakeTestServer(alphapointResponse)
	defer alphapointServer.Close()

	var witResponse wit.WitResponse
	json.Unmarshal(rawWitResponse, &witResponse)

	alphapoint.BaseUrl = alphapointServer.URL

	action := suite.Action
	isCreated := action.CreateNewExpense(witResponse, mocks.TestUserId)
	assert.True(suite.T(), isCreated)
}

func (suite *ActionSuite) TestCreatesNewExpenseFromCache() {
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
	alphapointResponse := `{
		"Realtime Currency Exchange Rate": {
			"5. Exchange Rate": ".7"
		}
	}`
	alphapointServer := mocks.MakeTestServer(alphapointResponse)

	var witResponse wit.WitResponse
	json.Unmarshal(rawWitResponse, &witResponse)

	alphapoint.BaseUrl = alphapointServer.URL

	action := suite.Action
	// Initial call without cache
	action.CreateNewExpense(witResponse, mocks.TestUserId)

	// Second call should not hit server
	alphapointServer.Close()
	isCreated := action.CreateNewExpense(witResponse, mocks.TestUserId)
	assert.True(suite.T(), isCreated)
}

func (suite *ActionSuite) TestGetsExpenseTotal() {
	action := suite.Action
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

	total, err := action.GetExpenseTotal(witResponse, mocks.TestUserId)
	assert.Equal(suite.T(), "0.00", total)
	assert.NoError(suite.T(), err)
}

func (suite *ActionSuite) TestReturnsErrorForInvalidPeriod() {
	action := &Actions{
		suite.Env.Db,
		suite.Env.Cache,
		suite.Env.Config,
	}

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
	_, err := action.GetExpenseTotal(witResponse, mocks.TestUserId)
	assert.EqualError(suite.T(), err, "foo is an invalid period")
}

func (suite *ActionSuite) TestCreatesNewUser() {
	action := suite.Action
	password := "my-password"

	isCreated := action.CreateNewUser(mocks.TestUserId, password)
	assert.True(suite.T(), isCreated)
}

func TestActionSuite(t *testing.T) {
	suite.Run(t, new(ActionSuite))
}
