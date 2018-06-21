// Package mocks provides testing utils and mocks to share across test suites.
package mocks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

// TestUserID is the common Telegram User ID we use in most test cases.
const TestUserID = uint(12345)

// TelegramMock mocks Telegram package.
type TelegramMock struct {
	Calls struct {
		SetWebhook int
		Send       int
		SendAction int
	}
}

// SessionMock mocks Session package.
type SessionMock struct {
	Calls struct {
		Get    int
		Set    int
		Delete int
	}
}

// Set mocks Session Set.
func (s *SessionMock) Set(cacheKey string, v interface{}, timeInSeconds int) {
	s.Calls.Set++
}

// Delete mocks Session Delete.
func (s *SessionMock) Delete(cacheKey string) error {
	s.Calls.Delete++
	return nil
}

// Get mocks Session Get.
func (s *SessionMock) Get(cacheKey string, v interface{}) error {
	s.Calls.Get++
	return nil
}

// SetWebhook mocks Telegram SetWebhook.
func (t *TelegramMock) SetWebhook() int {
	t.Calls.SetWebhook++
	statusCode := 200
	return statusCode
}

// Send mocks Telegram Send.
func (t *TelegramMock) Send(chatID int, message string) int {
	t.Calls.Send++
	statusCode := 200
	return statusCode
}

// SendAction mocks Telegram SendAction.
func (t *TelegramMock) SendAction(chatID int, action string) int {
	t.Calls.SendAction++
	statusCode := 200
	return statusCode
}

// MakeTestServer returns a test server with expected response.
func MakeTestServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))
}

// GetMockMessage returns a stub Telegram IncomingMessage.
func GetMockMessage(userResponse string) []byte {
	response := "Hello world"
	if userResponse != "" {
		response = userResponse
	}

	messageStr := fmt.Sprintf(`{
		"update_id": 123,
		"message": {
			"message_id": 123,
			"date": 20180314,
			"text": "%s",
			"from": {
				"id": 12345,
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
	}`, response)
	message := []byte(messageStr)
	return message
}

// MockTime implements CurrentTime interface.
type MockTime struct {
	CurrentTime time.Time
}

// Now returns the current time.
func (m *MockTime) Now() time.Time {
	return m.CurrentTime
}

// MessageMapMock replaces conversatio MessageMap with predictable results.
var MessageMapMock = map[string][]string{
	"default": []string{
		"This is a default message",
	},
	"get_expense_total_invalid_period": []string{
		"Whoops!",
	},
	"get_expense_total_error": []string{
		"Whoops!",
	},
	"get_expense_total_success": []string{
		"You spent {{var}}",
	},
	"get_expense_total_password_invalid": []string{
		"This password is invalid",
	},
	"get_expense_total_ask_for_password": []string{
		"I need your password",
	},
	"track_expense_error": []string{
		"Whoops!",
	},
	"track_expense_success": []string{
		"Roger that!",
	},
	"onboard_user_ask_for_password": []string{
		"What's your password?",
	},
	"onboard_user_confirm_password_error": []string{
		"I didn't understand that",
	},
	"onboard_user_reject_password": []string{
		"Okay try again later",
	},
	"onboard_user_confirm_password": []string{
		"Your password is {{var}}",
	},
	"onboard_user_say_outro": []string{
		"Outro message",
	},
	"onboard_user_account_creation_failed": []string{
		"Couldn't create account",
	},
}
