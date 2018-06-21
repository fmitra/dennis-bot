package telegram

import (
	"fmt"
	"testing"

	mocks "github.com/fmitra/dennis-bot/test"
	"github.com/stretchr/testify/assert"
)

func TestTelegram(t *testing.T) {
	t.Run("Returns client with default config", func(t *testing.T) {
		telegram := NewClient("telegramToken", "https://localhost")

		assert.Equal(t, BaseURL, telegram.BaseURL)
	})

	t.Run("Sets webhook", func(t *testing.T) {
		server := mocks.MakeTestServer("")
		defer server.Close()
		telegram := &Client{
			Token:   "telegramToken",
			Domain:  "https://localhost",
			BaseURL: fmt.Sprintf("%s/", server.URL),
		}

		statusCode := telegram.SetWebhook()
		assert.Equal(t, 200, statusCode)
	})

	t.Run("Sends telegram message", func(t *testing.T) {
		server := mocks.MakeTestServer("")
		defer server.Close()
		telegram := &Client{
			Token:   "telegramToken",
			Domain:  "https://localhost",
			BaseURL: fmt.Sprintf("%s/", server.URL),
		}

		chatID := 5
		message := "Hello world"
		statusCode := telegram.Send(chatID, message)
		assert.Equal(t, 200, statusCode)
	})

	t.Run("Sends a telegram chat action", func(t *testing.T) {
		server := mocks.MakeTestServer("")
		defer server.Close()
		telegram := &Client{
			Token:   "telegramToken",
			Domain:  "https://localhost",
			BaseURL: fmt.Sprintf("%s/", server.URL),
		}

		chatID := 5
		action := "typing"
		statusCode := telegram.SendAction(chatID, action)
		assert.Equal(t, 200, statusCode)
	})
}
