package telegram

import (
	"testing"

	mocks "github.com/fmitra/dennis-bot/test"
	"github.com/stretchr/testify/assert"
)

var incomingMessage = IncomingMessage{
	UpdateID: 1,
	Message: Message{
		MessageID: 2,
		Date:      1521367820,
		Text:      "Hello world",
		From: User{
			ID:        mocks.TestUserID,
			FirstName: "Jane",
			LastName:  "Doe",
			UserName:  "janedoe123",
		},
		Chat: Chat{
			ID:        3,
			FirstName: "John",
			LastName:  "Doe",
			UserName:  "johndoe456",
		},
	},
}

func TestMessages(t *testing.T) {
	t.Run("Gets Client ID", func(t *testing.T) {
		assert.Equal(t, 3, incomingMessage.GetChatID())
	})

	t.Run("Gets User", func(t *testing.T) {
		user := User{
			ID:        mocks.TestUserID,
			FirstName: "Jane",
			LastName:  "Doe",
			UserName:  "janedoe123",
		}
		assert.Equal(t, user, incomingMessage.GetUser())
	})

	t.Run("Gets Message", func(t *testing.T) {
		assert.Equal(t, "Hello world", incomingMessage.GetMessage())
	})
}
