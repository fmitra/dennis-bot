package internal

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis-bot/config"
	"github.com/fmitra/dennis-bot/pkg/alphapoint"
	"github.com/fmitra/dennis-bot/pkg/expenses"
	"github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/wit"
	mocks "github.com/fmitra/dennis-bot/test"
)

func TestBot(t *testing.T) {
	t.Run("Should handle tracking intent from Wit.ai", func(t *testing.T) {
		var incMessage telegram.IncomingMessage
		message := mocks.GetMockMessage()
		json.Unmarshal(message, &incMessage)
		MessageMap = mocks.MessageMapMock

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
		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}

		response := bot.BuildResponse(incMessage)
		assert.Equal(t, BotResponse("Roger that!"), response)
	})

	t.Run("Should handle period query intent from Wit.ai", func(t *testing.T) {
		var incMessage telegram.IncomingMessage
		message := mocks.GetMockMessage()
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

		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}
		bot.env.db.Where("user_id = ?", mocks.TestUserId).
			Unscoped().
			Delete(expenses.Expense{})

		response := bot.BuildResponse(incMessage)
		assert.Equal(t, BotResponse("You spent 0.00"), response)
	})

	t.Run("Should return a default message", func(t *testing.T) {
		var incMessage telegram.IncomingMessage
		message := mocks.GetMockMessage()
		json.Unmarshal(message, &incMessage)
		MessageMap = mocks.MessageMapMock

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

		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}

		response := bot.BuildResponse(incMessage)
		assert.Equal(t, BotResponse("This is a default message"), response)
	})

	t.Run("Should return a error message", func(t *testing.T) {
		var incMessage telegram.IncomingMessage
		message := mocks.GetMockMessage()
		json.Unmarshal(message, &incMessage)
		MessageMap = mocks.MessageMapMock

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

		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}

		response := bot.BuildResponse(incMessage)
		assert.Equal(t, BotResponse("Whoops!"), response)
	})

	t.Run("Receives and responds through telegram", func(t *testing.T) {
		configFile := "../config/config.json"
		telegramMock := &mocks.TelegramMock{}
		sessionMock := &mocks.SessionMock{}
		env := &Env{
			cache:    sessionMock,
			config:   config.LoadConfig(configFile),
			telegram: telegramMock,
		}

		bot := &Bot{env}
		message := mocks.GetMockMessage()
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
		assert.Equal(t, 1, telegramMock.Calls.Send)
		assert.Equal(t, 1, sessionMock.Calls.Set)
	})

	t.Run("Receives an incoming message", func(t *testing.T) {
		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}
		message := mocks.GetMockMessage()

		var incMessage telegram.IncomingMessage
		json.Unmarshal(message, &incMessage)

		receivedMessage, _ := bot.ReceiveMessage(message)

		assert.Equal(t, incMessage, receivedMessage)
	})

	t.Run("Sends an outgoing message", func(t *testing.T) {
		telegramServer := mocks.MakeTestServer("")
		telegram.BaseUrl = fmt.Sprintf("%s/", telegramServer.URL)
		defer telegramServer.Close()

		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}
		message := mocks.GetMockMessage()
		response := BotResponse("Hello world")

		var incMessage telegram.IncomingMessage
		json.Unmarshal(message, &incMessage)

		statusCode := bot.SendMessage(response, incMessage)
		assert.Equal(t, 200, statusCode)
	})

}
