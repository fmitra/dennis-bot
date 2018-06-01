package telegram

import (
	"fmt"
	"testing"

	"github.com/fmitra/dennis/mocks"
	"github.com/stretchr/testify/assert"
)

func TestTelegram(t *testing.T) {
	t.Run("Returns client with default config", func(t *testing.T) {
		telegram := NewClient("telegramToken", "https://localhost")

		assert.Equal(t, BaseUrl, telegram.BaseUrl)
	})

	t.Run("Sets webhook", func(t *testing.T) {
		server := mocks.MakeTestServer("")
		telegram := Client{
			Token:   "telegramToken",
			Domain:  "https://localhost",
			BaseUrl: fmt.Sprintf("%s/", server.URL),
		}

		statusCode := telegram.SetWebhook()
		assert.Equal(t, 200, statusCode)
	})

	t.Run("Sends telegram message", func(t *testing.T) {
		server := mocks.MakeTestServer("")
		telegram := Client{
			Token:   "telegramToken",
			Domain:  "https://localhost",
			BaseUrl: fmt.Sprintf("%s/", server.URL),
		}

		chatId := 5
		message := "Hello world"
		statusCode := telegram.Send(chatId, message)
		assert.Equal(t, 200, statusCode)
	})
}
