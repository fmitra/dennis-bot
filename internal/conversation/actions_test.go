package conversation

import (
	"crypto/rsa"
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

	var witResponse wit.Response
	json.Unmarshal(rawWitResponse, &witResponse)

	alphapoint.BaseURL = alphapointServer.URL

	action := suite.Action
	publicKey := rsa.PublicKey{}
	err := action.CreateNewExpense(witResponse, mocks.TestUserID, publicKey)
	assert.NoError(suite.T(), err)
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

	var witResponse wit.Response
	json.Unmarshal(rawWitResponse, &witResponse)

	alphapoint.BaseURL = alphapointServer.URL

	publicKey := rsa.PublicKey{}

	action := suite.Action
	// Initial call without cache
	action.CreateNewExpense(witResponse, mocks.TestUserID, publicKey)

	// Second call should not hit server
	alphapointServer.Close()
	err := action.CreateNewExpense(witResponse, mocks.TestUserID, publicKey)
	assert.NoError(suite.T(), err)
}

func (suite *ActionSuite) TestGetsExpenseTotal() {
	action := suite.Action
	privateKey := rsa.PrivateKey{}
	period := "month"
	total, err := action.GetExpenseTotal(period, uint(200), privateKey)
	assert.Equal(suite.T(), "0.00", total)
	assert.NoError(suite.T(), err)
}

func (suite *ActionSuite) TestReturnsErrorForInvalidPeriod() {
	action := &Actions{
		suite.Env.Db,
		suite.Env.Cache,
		suite.Env.Config,
	}

	privateKey := rsa.PrivateKey{}
	period := "foo"
	_, err := action.GetExpenseTotal(period, mocks.TestUserID, privateKey)
	assert.EqualError(suite.T(), err, "foo is an invalid period")
}

func (suite *ActionSuite) TestCreatesNewUser() {
	action := suite.Action
	password := "my-password"

	err := action.CreateNewUser(mocks.TestUserID, password)
	assert.NoError(suite.T(), err)
}

func TestActionSuite(t *testing.T) {
	suite.Run(t, new(ActionSuite))
}
