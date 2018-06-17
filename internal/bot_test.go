package internal

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	convo "github.com/fmitra/dennis-bot/internal/conversation"
	"github.com/fmitra/dennis-bot/pkg/alphapoint"
	"github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/users"
	"github.com/fmitra/dennis-bot/pkg/wit"
	mocks "github.com/fmitra/dennis-bot/test"
)

type BotSuite struct {
	suite.Suite
	Env    *mocks.TestEnv
	BotEnv *Env
}

func (suite *BotSuite) SetupSuite() {
	configFile := "../config/config.json"
	testEnv := mocks.GetTestEnv(configFile)
	telegram := telegram.NewClient(
		testEnv.Config.Telegram.Token,
		testEnv.Config.BotDomain,
	)
	suite.Env = testEnv
	suite.BotEnv = &Env{
		suite.Env.Db,
		suite.Env.Cache,
		suite.Env.Config,
		telegram,
	}
}

func (suite *BotSuite) BeforeTest(suiteName, testName string) {
	convo.MessageMap = mocks.MessageMapMock
	mocks.CreateTestUser(suite.Env.Db)
}

func (suite *BotSuite) AfterTest(suiteName, testName string) {
	mocks.CleanUpEnv(suite.Env)
}

func GetTestUser(db *gorm.DB) users.User {
	var user users.User
	db.Where("telegram_id = ?", mocks.TestUserId).First(&user)
	return user
}

func (suite *BotSuite) TestHandlesTrackingIntent() {
	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("")
	json.Unmarshal(message, &incMessage)

	witResponse := `{
		"entities": {
			"amount": [
				{ "value": "20 USD", "confidence": 100.00 }
			],
			"datetime": [
				{ "value": "", "confidence": 100.00 }
			],
			"description": [
				{ "value": "Food", "confidence": 100.00 }
			]
		}
	}`
	alphapointResponse := `{
		"Realtime Currency Exchange Rate": {
			"5. Exchange Rate": ".7"
		}
	}`
	witServer := mocks.MakeTestServer(witResponse)
	alphapointServer := mocks.MakeTestServer(alphapointResponse)
	defer witServer.Close()
	defer alphapointServer.Close()

	wit.BaseUrl = witServer.URL
	alphapoint.BaseUrl = alphapointServer.URL
	bot := &Bot{suite.BotEnv}

	response := bot.BuildResponse(incMessage)
	assert.Equal(suite.T(), convo.BotResponse("Roger that!"), response)
}

func (suite *BotSuite) TestHandlesExpenseTotalIntent() {
	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("")
	json.Unmarshal(message, &incMessage)

	witResponse := `{
		"entities": {
			"amount": [],
			"datetime": [],
			"description": [],
			"total_spent": [
				{ "value": "month", "confidence": 100.00 }
			]
		}
	}`
	witServer := mocks.MakeTestServer(witResponse)
	wit.BaseUrl = witServer.URL
	defer witServer.Close()

	bot := &Bot{suite.BotEnv}

	response := bot.BuildResponse(incMessage)
	assert.Equal(suite.T(), convo.BotResponse("I need your password"), response)
}

func (suite *BotSuite) TestReturnsDefaultMessage() {
	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("")
	json.Unmarshal(message, &incMessage)

	witResponse := `{
		"entities": {
			"amount": [],
			"datetime": [],
			"description": [],
			"total_spent": []
		}
	}`
	witServer := mocks.MakeTestServer(witResponse)
	wit.BaseUrl = witServer.URL
	defer witServer.Close()

	bot := &Bot{suite.BotEnv}

	response := bot.BuildResponse(incMessage)
	assert.Equal(suite.T(), convo.BotResponse("This is a default message"), response)
}

func (suite *BotSuite) TestReturnsErrorMessage() {
	var incMessage telegram.IncomingMessage
	message := mocks.GetMockMessage("")
	json.Unmarshal(message, &incMessage)

	witResponse := `{
		"entities": {
			"amount": [],
			"datetime": [],
			"description": [],
			"total_spent": [
				{ "value": "foo", "confidence": 100.00 }
			]
		}
	}`
	witServer := mocks.MakeTestServer(witResponse)
	wit.BaseUrl = witServer.URL
	defer witServer.Close()

	bot := &Bot{suite.BotEnv}

	cacheKey := fmt.Sprintf("%s_password", strconv.Itoa(int(mocks.TestUserId)))
	suite.Env.Cache.Set(cacheKey, "my-password", 180)

	response := bot.BuildResponse(incMessage)
	assert.Equal(suite.T(), convo.BotResponse("Whoops!"), response)
}

func (suite *BotSuite) TestReceivesRespondsWithTelegram() {
	telegramMock := &mocks.TelegramMock{}
	sessionMock := &mocks.SessionMock{}

	bot := &Bot{
		&Env{
			suite.Env.Db,
			sessionMock,
			suite.Env.Config,
			telegramMock,
		},
	}

	message := mocks.GetMockMessage("")
	witResponse := `{
		"entities": {
			"amount": [],
			"datetime": [],
			"description": []
		}
	}`
	witServer := mocks.MakeTestServer(witResponse)
	telegramServer := mocks.MakeTestServer("")
	wit.BaseUrl = witServer.URL
	telegram.BaseUrl = fmt.Sprintf("%s/", telegramServer.URL)
	defer witServer.Close()
	defer telegramServer.Close()

	bot.Converse(message)
	assert.Equal(suite.T(), 1, telegramMock.Calls.Send)
}

func (suite *BotSuite) TestReceivesIncomingMessage() {
	bot := &Bot{suite.BotEnv}
	message := mocks.GetMockMessage("")

	var incMessage telegram.IncomingMessage
	json.Unmarshal(message, &incMessage)

	receivedMessage, _ := bot.ReceiveMessage(message)

	assert.Equal(suite.T(), incMessage, receivedMessage)
}

func (suite *BotSuite) TestSendsOutgoingMessage() {
	telegramServer := mocks.MakeTestServer("")
	telegram.BaseUrl = fmt.Sprintf("%s/", telegramServer.URL)
	telegramMock := telegram.NewClient("", "")
	defer telegramServer.Close()

	bot := &Bot{
		&Env{
			suite.Env.Db,
			suite.Env.Cache,
			suite.Env.Config,
			telegramMock,
		},
	}
	message := mocks.GetMockMessage("")
	response := convo.BotResponse("Hello world")

	var incMessage telegram.IncomingMessage
	json.Unmarshal(message, &incMessage)

	statusCode := bot.SendMessage(response, incMessage)
	assert.Equal(suite.T(), 200, statusCode)
}

func TestBotSuite(t *testing.T) {
	suite.Run(t, new(BotSuite))
}
