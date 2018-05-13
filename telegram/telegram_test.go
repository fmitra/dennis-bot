package telegram

import (
	"testing"
	"net/http"
	"io"
	"io/ioutil"
	"bytes"
)

type HttpMock struct {
	Calls struct {
		Post int
		Get int
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

		if mockHttpLib.Calls.Get != 1 {
			t.Error("Http GET request failed")
		}
	})

	t.Run("Sets webhook on init", func(t *testing.T) {
		token := "telegramToken"
		domain := "https://localhost"
		mockHttpLib := &HttpMock{}
		Init(token, domain, mockHttpLib)

		if Client.Token != "telegramToken" {
			t.Error("Client not initialized with token")
		}

		if Client.Domain != "https://localhost" {
			t.Error("Client not initialized with domain")
		}

		if Client.Http != mockHttpLib {
			t.Error("Client not initialized with Http library")
		}
	})

	t.Run("Sets webhook", func(t *testing.T) {
		mock := &HttpMock{}
		telegram := Telegram{
			Token: "telegramToken",
			Domain: "https://localhost",
			BaseUrl: "https://api.telegram.org/bot",
			Http: mock,
		}

		telegram.SetWebhook()
		if mock.Calls.Get != 1 {
			t.Error("Http GET request failed")
		}
	})

	t.Run("Sends telegram message", func(t *testing.T) {
		mock := &HttpMock{}
		telegram := Telegram{
			Token: "telegramToken",
			Domain: "https://localhost",
			BaseUrl: "https://api.telegram.org/bot",
			Http: mock,
		}

		chatId := 5
		message := "Hello world"
		telegram.Send(chatId, message)
		if mock.Calls.Post != 1 {
			t.Error("Http POST request failed")
		}
	})
}
