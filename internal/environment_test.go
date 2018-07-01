package internal

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/fmitra/dennis-bot/pkg/alphapoint"
	"github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/wit"
	mocks "github.com/fmitra/dennis-bot/test"
)

type EnvSuite struct {
	suite.Suite
	Env *mocks.TestEnv
}

func (suite *EnvSuite) SetupSuite() {
	configFile := "../config/config.json"
	suite.Env = mocks.GetTestEnv(configFile)
}

func (suite *EnvSuite) TestRespondsToHealthCheck() {
	env := &Env{
		suite.Env.Db,
		suite.Env.Cache,
		suite.Env.Config,
		&telegram.Client{},
		&wit.Client{},
		&alphapoint.Client{},
	}

	req, err := http.NewRequest("GET", "/healthcheck", nil)
	assert.NoError(suite.T(), err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.HealthCheck())

	handler.ServeHTTP(rr, req)
	assert.Equal(suite.T(), "ok", rr.Body.String())
}

func (suite *EnvSuite) TestRespondsToWebhook() {
	message := mocks.GetMockMessage("")
	witResponse := `{
		"entities": {
			"amount": [],
			"datetime": [],
			"description": []
		}
	}`
	telegramServer := mocks.MakeTestServer("")
	witServer := mocks.MakeTestServer(witResponse)

	alphapointClient := alphapoint.NewClient("")
	witClient := wit.NewClient("")
	telegramClient := telegram.NewClient("", "")
	telegramClient.BaseURL = fmt.Sprintf("%s/", telegramServer.URL)
	witClient.BaseURL = witServer.URL

	env := &Env{
		suite.Env.Db,
		suite.Env.Cache,
		suite.Env.Config,
		telegramClient,
		witClient,
		alphapointClient,
	}

	req, err := http.NewRequest("POST", "/webook", bytes.NewBuffer(message))
	assert.NoError(suite.T(), err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.Webhook())

	handler.ServeHTTP(rr, req)
	assert.Equal(suite.T(), "received", rr.Body.String())

	// Business logic is handled in a go routine, so we need
	// to add a delay to close the test server
	time.Sleep(time.Second * 1)
	telegramServer.Close()
	witServer.Close()
}

func (suite *EnvSuite) TestShouldLoadFromConfig() {
	env := LoadEnv(suite.Env.Config)

	assert.NotNil(suite.T(), env.db)
	assert.NotNil(suite.T(), env.cache)
	assert.NotNil(suite.T(), env.config)
	assert.NotNil(suite.T(), env.telegram)
}

func TestEnvSuite(t *testing.T) {
	suite.Run(t, new(EnvSuite))
}
