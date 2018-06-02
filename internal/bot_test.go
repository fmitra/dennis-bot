package internal

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis/alphapoint"
	"github.com/fmitra/dennis/config"
	"github.com/fmitra/dennis/expenses"
	"github.com/fmitra/dennis/mocks"
	"github.com/fmitra/dennis/telegram"
	"github.com/fmitra/dennis/wit"
)

func TestBot(t *testing.T) {
	t.Run("It should generate a message with optional message var", func(t *testing.T) {
		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}
		var message string
		MessageMap = mocks.MessageMapMock

		message = bot.GetMessage("tracking_success", "")
		assert.Equal(t, "Roger that!", message)

		message = bot.GetMessage("period_total_success", "20")
		assert.Equal(t, "You spent 20", message)
	})

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
		assert.Equal(t, "Roger that!", response)
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
		assert.Equal(t, "You spent 0.00", response)
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
		assert.Equal(t, "This is a default message", response)
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
		assert.Equal(t, "Whoops!", response)
	})

	t.Run("It should return a historical total by period", func(t *testing.T) {
		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}
		bot.env.db.Where("user_id = ?", mocks.TestUserId).
			Unscoped().
			Delete(expenses.Expense{})

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

		strTotal, err := bot.GetTotalByPeriod(witResponse, mocks.TestUserId)
		assert.Equal(t, "0.00", strTotal)
		assert.NoError(t, err)
	})

	t.Run("It should return error for invalid period", func(t *testing.T) {
		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}
		bot.env.db.Where("user_id = ?", mocks.TestUserId).
			Unscoped().
			Delete(expenses.Expense{})

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

		strTotal, err := bot.GetTotalByPeriod(witResponse, mocks.TestUserId)
		assert.Equal(t, "0.00", strTotal)
		assert.EqualError(t, err, "foo is an invalid period")
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
		response := "Hello world"

		var incMessage telegram.IncomingMessage
		json.Unmarshal(message, &incMessage)

		statusCode := bot.SendMessage(response, incMessage)
		assert.Equal(t, 200, statusCode)
	})

	t.Run("Creates a new expense", func(t *testing.T) {
		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		bot := &Bot{env}
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
		userId := mocks.TestUserId

		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)

		alphapoint.BaseUrl = alphapointServer.URL

		isCreated := bot.NewExpense(witResponse, userId)
		assert.True(t, isCreated)
	})

	t.Run("Creates a new expense from cache", func(t *testing.T) {
		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		env.cache.Delete("SGD_USD")
		bot := &Bot{env}
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
		userId := mocks.TestUserId

		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)

		alphapoint.BaseUrl = alphapointServer.URL

		// Initial call without cache
		bot.NewExpense(witResponse, userId)

		// Second call should not hit server
		alphapointServer.Close()
		isCreated := bot.NewExpense(witResponse, userId)
		assert.True(t, isCreated)
	})
}
