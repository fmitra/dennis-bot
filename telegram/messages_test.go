package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var incomingMessage = IncomingMessage{
	UpdateId: 1,
	Message: Message{
		MessageId: 2,
		Date:      1521367820,
		Text:      "Hello world",
		From: User{
			Id:        123,
			FirstName: "Jane",
			LastName:  "Doe",
			UserName:  "janedoe123",
		},
		Chat: Chat{
			Id:        3,
			FirstName: "John",
			LastName:  "Doe",
			UserName:  "johndoe456",
		},
	},
}

func TestMessages(t *testing.T) {
	t.Run("Gets Client ID", func(t *testing.T) {
		assert.Equal(t, 3, incomingMessage.GetChatId())
	})

	t.Run("Gets User", func(t *testing.T) {
		user := User{
			Id:        123,
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
