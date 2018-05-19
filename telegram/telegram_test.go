package telegram

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeTestServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))
}

func TestTelegram(t *testing.T) {
	t.Run("Returns client with default config", func(t *testing.T) {
		telegram := Client("telegramToken", "https://localhost")

		assert.Equal(t, BaseUrl, telegram.BaseUrl)
	})

	t.Run("Sets webhook", func(t *testing.T) {
		server := makeTestServer("")
		telegram := client{
			Token:   "telegramToken",
			Domain:  "https://localhost",
			BaseUrl: fmt.Sprintf("%s/", server.URL),
		}

		statusCode := telegram.SetWebhook()
		assert.Equal(t, 200, statusCode)
	})

	t.Run("Sends telegram message", func(t *testing.T) {
		server := makeTestServer("")
		telegram := client{
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
