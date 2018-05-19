package mocks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

type TelegramMock struct {
	Calls struct {
		SetWebhook int
		Send       int
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
	statusCode := 200
	return statusCode
}

func (t *TelegramMock) Send(chatId int, message string) int {
	t.Calls.Send++
	statusCode := 200
	return statusCode
}

func MakeTestServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))
}

func GetMockMessage() []byte {
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
