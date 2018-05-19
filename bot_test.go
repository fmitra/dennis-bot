package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis/alphapoint"
	"github.com/fmitra/dennis/config"
	"github.com/fmitra/dennis/mocks"
	"github.com/fmitra/dennis/telegram"
	"github.com/fmitra/dennis/wit"
)

func TestBot(t *testing.T) {
	t.Run("Receives and responds through telegram", func(t *testing.T) {
		configFile := "config/config.json"
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

		<-bot.Converse(message)
		assert.Equal(t, 1, telegramMock.Calls.Send)
		assert.Equal(t, 1, sessionMock.Calls.Set)
	})

	t.Run("Receives an incoming message", func(t *testing.T) {
		configFile := "config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}
		message := mocks.GetMockMessage()

		var incMessage telegram.IncomingMessage
		json.Unmarshal(message, &incMessage)

		receivedMessage, _ := bot.ReceiveMessage(message)

		assert.Equal(t, incMessage, receivedMessage)
	})

	t.Run("Sends an outgoing message", func(t *testing.T) {
		configFile := "config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}
		message := mocks.GetMockMessage()
		keyword := "track.success"

		var incMessage telegram.IncomingMessage
		json.Unmarshal(message, &incMessage)

		server := mocks.MakeTestServer("")

		// We need to add a trailing slash because the telegram token
		// format will be treated as an invalid port on the test server
		telegram.BaseUrl = fmt.Sprintf("%s/", server.URL)

		statusCode := bot.SendMessage(keyword, incMessage)
		assert.Equal(t, 200, statusCode)
	})

	t.Run("Retrieves a string response to send", func(t *testing.T) {
		configFile := "config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}

		// messageMap var is setup with multiple possible responses.
		// We will instead replace it to only offer one response for
		// simplicity
		messageMap = map[string][]string{
			"default": []string{
				"This is a default message",
			},
			"track.success": []string{
				"This is a successful tracking message",
			},
			"track.error": []string{
				"This is a failed tracking message",
			},
		}

		successMessage := bot.GetResponse("track.success")
		errorMessage := bot.GetResponse("track.error")
		defaultMessage := bot.GetResponse("default")

		assert.Equal(t, defaultMessage, "This is a default message")
		assert.Equal(t, successMessage, "This is a successful tracking message")
		assert.Equal(t, errorMessage, "This is a failed tracking message")
	})

	t.Run("Maps incoming message to tracking keyword", func(t *testing.T) {
		configFile := "config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}
		message := mocks.GetMockMessage()
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

		wit.BaseUrl = witServer.URL
		alphapoint.BaseUrl = alphapointServer.URL

		var incMessage telegram.IncomingMessage
		json.Unmarshal(message, &incMessage)

		keyword := bot.MapToKeyword(incMessage)

		assert.Equal(t, "track.success", keyword)
	})

	t.Run("Maps incoming message to tracking error keyword", func(t *testing.T) {
		configFile := "config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}
		message := mocks.GetMockMessage()
		witResponse := `{
			"entities": {
				"amount": [
					{ "value": "20 USD", "confidence": 100.00 }
				],
				"datetime": [
					{ "value": "", "confidence": 100.00 }
				],
				"description": []
			}
		}`
		alphapointResponse := `{
			"Realtime Currency Exchange Rate": {
				"5. Exchange Rate": ".7"
			}
		}`
		witServer := mocks.MakeTestServer(witResponse)
		alphapointServer := mocks.MakeTestServer(alphapointResponse)

		wit.BaseUrl = witServer.URL
		alphapoint.BaseUrl = alphapointServer.URL

		var incMessage telegram.IncomingMessage
		json.Unmarshal(message, &incMessage)

		keyword := bot.MapToKeyword(incMessage)

		assert.Equal(t, "track.error", keyword)
	})

	t.Run("Maps incoming message to default keyword", func(t *testing.T) {
		configFile := "config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}
		message := mocks.GetMockMessage()
		witResponse := `{
			"entities": {
				"amount": [],
				"datetime": [],
				"description": []
			}
		}`
		alphapointResponse := `{
			"Realtime Currency Exchange Rate": {
				"5. Exchange Rate": ".7"
			}
		}`
		witServer := mocks.MakeTestServer(witResponse)
		alphapointServer := mocks.MakeTestServer(alphapointResponse)

		wit.BaseUrl = witServer.URL
		alphapoint.BaseUrl = alphapointServer.URL

		var incMessage telegram.IncomingMessage
		json.Unmarshal(message, &incMessage)

		keyword := bot.MapToKeyword(incMessage)

		assert.Equal(t, "default", keyword)
	})

	t.Run("Creates a new expense", func(t *testing.T) {
		configFile := "config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}
		rawWitResponse := []byte(`{
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
		}`)
		alphapointResponse := `{
			"Realtime Currency Exchange Rate": {
				"5. Exchange Rate": ".7"
			}
		}`
		alphapointServer := mocks.MakeTestServer(alphapointResponse)
		userId := 123

		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)

		alphapoint.BaseUrl = alphapointServer.URL

		bot.NewExpense(witResponse, userId)
	})
}
