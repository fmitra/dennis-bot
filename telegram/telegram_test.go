package telegram

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HttpMock struct {
	Calls struct {
		Post int
		Get  int
	}
}

func (h *HttpMock) Get(url string) (*http.Response, error) {
	h.Calls.Get++
	response := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer([]byte("GET"))),
	}
	return response, nil
}

func (h *HttpMock) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	h.Calls.Post++
	response := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer([]byte("POST"))),
	}
	return response, nil
}

func TestTelegram(t *testing.T) {
	t.Run("Sets Client on init", func(t *testing.T) {
		token := "telegramToken"
		domain := "https://localhost"
		mockHttpLib := &HttpMock{}

		<-Init(token, domain, mockHttpLib)

		assert.Equal(t, 1, mockHttpLib.Calls.Get)
	})

	t.Run("Sets webhook on init", func(t *testing.T) {
		token := "telegramToken"
		domain := "https://localhost"
		mockHttpLib := &HttpMock{}
		Init(token, domain, mockHttpLib)

		assert.Equal(t, "telegramToken", Client.Token)
		assert.Equal(t, "https://localhost", Client.Domain)
		assert.Equal(t, mockHttpLib, Client.Http)
	})

	t.Run("Sets webhook", func(t *testing.T) {
		mock := &HttpMock{}
		telegram := Telegram{
			Token:   "telegramToken",
			Domain:  "https://localhost",
			BaseUrl: "https://api.telegram.org/bot",
			Http:    mock,
		}

		telegram.SetWebhook()
		assert.Equal(t, 1, mock.Calls.Get)
	})

	t.Run("Sends telegram message", func(t *testing.T) {
		mock := &HttpMock{}
		telegram := Telegram{
			Token:   "telegramToken",
			Domain:  "https://localhost",
			BaseUrl: "https://api.telegram.org/bot",
			Http:    mock,
		}

		chatId := 5
		message := "Hello world"
		telegram.Send(chatId, message)
		assert.Equal(t, 1, mock.Calls.Post)
	})
}
