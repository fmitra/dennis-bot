package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis/alphapoint"
	"github.com/fmitra/dennis/config"
	"github.com/fmitra/dennis/telegram"
	"github.com/fmitra/dennis/wit"
)

type TelegramMock struct {
	Calls struct {
		SetWebhook int
		Send int
	}
}

type SessionMock struct {
	Calls struct {
		Get int
		Set int
	}
}

func (s *SessionMock) Set(cacheKey string, v interface{}) {
	s.Calls.Set++
}

func (s *SessionMock) Get(cacheKey string, v interface{}) error {
	s.Calls.Get++
	return nil
}

func (t *TelegramMock) SetWebhook() int {
	t.Calls.SetWebhook++
	return int(200)
}

func (t *TelegramMock) Send(chatId int, message string) int {
	t.Calls.Send++
	return int(200)
}

func makeTestServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))
}

func getMockMessage() []byte {
	message := []byte(`{
		"update_id": 123,
		"message": {
			"message_id": 123,
			"date": 20180314,
			"text": "Hello world",
			"from": {
				"id": 456,
				"first_name": "Jane",
				"last_name": "Doe",
				"username": "janedoe"
			},
			"chat": {
				"id": 456,
				"first_name": "Jane",
				"last_name": "Doe",
				"username": "janedoe"
			}
		}
	}`)
	return message
}

func TestBot(t *testing.T) {
	t.Run("Receives and responds through telegram", func(t *testing.T) {
		telegramMock := &TelegramMock{}
		sessionMock := &SessionMock{}
		env := &Env{
			cache: sessionMock,
			config: config.LoadConfig(),
			telegram: telegramMock,
		}

		bot := &Bot{env}
		message := getMockMessage()
		witResponse := `{
			"entities": {
				"amount": [],
				"datetime": [],
				"description": []
			}
		}`
		witServer := makeTestServer(witResponse)
		telegramServer := makeTestServer("")
		wit.BaseUrl = witServer.URL
		telegram.BaseUrl = fmt.Sprintf("%s/", telegramServer.URL)

		<- bot.Converse(message)
		assert.Equal(t, 1, telegramMock.Calls.Send)
		assert.Equal(t, 1, sessionMock.Calls.Set)
	})

	t.Run("Receives an incoming message", func(t *testing.T) {
		env := LoadEnv(config.LoadConfig())
		bot := &Bot{env}
		message := getMockMessage()

		var incMessage telegram.IncomingMessage
		json.Unmarshal(message, &incMessage)

		receivedMessage, _ := bot.ReceiveMessage(message)

		assert.Equal(t, incMessage, receivedMessage)
	})

	t.Run("Sends an outgoing message", func(t *testing.T) {
		env := LoadEnv(config.LoadConfig())
		bot := &Bot{env}
		message := getMockMessage()
		keyword := "track.success"

		var incMessage telegram.IncomingMessage
		json.Unmarshal(message, &incMessage)

		server := makeTestServer("")

		// We need to add a trailing slash because the telegram token
		// format will be treated as an invalid port on the test server
		telegram.BaseUrl = fmt.Sprintf("%s/", server.URL)

		statusCode := bot.SendMessage(keyword, incMessage)
		assert.Equal(t, 200, statusCode)
	})

	t.Run("Retrieves a string response to send", func(t *testing.T) {
		env := LoadEnv(config.LoadConfig())
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
		env := LoadEnv(config.LoadConfig())
		bot := &Bot{env}
		message := getMockMessage()
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
		witServer := makeTestServer(witResponse)
		alphapointServer := makeTestServer(alphapointResponse)

		wit.BaseUrl = witServer.URL
		alphapoint.BaseUrl = alphapointServer.URL

		var incMessage telegram.IncomingMessage
		json.Unmarshal(message, &incMessage)

		keyword := bot.MapToKeyword(incMessage)

		assert.Equal(t, "track.success", keyword)
	})

	t.Run("Maps incoming message to tracking error keyword", func(t *testing.T) {
		env := LoadEnv(config.LoadConfig())
		bot := &Bot{env}
		message := getMockMessage()
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
		witServer := makeTestServer(witResponse)
		alphapointServer := makeTestServer(alphapointResponse)

		wit.BaseUrl = witServer.URL
		alphapoint.BaseUrl = alphapointServer.URL

		var incMessage telegram.IncomingMessage
		json.Unmarshal(message, &incMessage)

		keyword := bot.MapToKeyword(incMessage)

		assert.Equal(t, "track.error", keyword)
	})

	t.Run("Maps incoming message to default keyword", func(t *testing.T) {
		env := LoadEnv(config.LoadConfig())
		bot := &Bot{env}
		message := getMockMessage()
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
		witServer := makeTestServer(witResponse)
		alphapointServer := makeTestServer(alphapointResponse)

		wit.BaseUrl = witServer.URL
		alphapoint.BaseUrl = alphapointServer.URL

		var incMessage telegram.IncomingMessage
		json.Unmarshal(message, &incMessage)

		keyword := bot.MapToKeyword(incMessage)

		assert.Equal(t, "default", keyword)
	})

	t.Run("Creates a new expense", func(t *testing.T) {
		env := LoadEnv(config.LoadConfig())
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
		alphapointServer := makeTestServer(alphapointResponse)
		userId := 123

		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)

		alphapoint.BaseUrl = alphapointServer.URL

		bot.NewExpense(witResponse, userId)
	})
}
