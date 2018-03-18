package telegram

import (
	"testing"
)

var incomingMessage = IncomingMessage{
	UpdateId: 1,
	Message: Message{
		MessageId: 2,
		Date: 1521367820,
		Text: "Hello world",
		From: User{
			Id: 123,
			FirstName: "Jane",
			LastName: "Doe",
			UserName: "janedoe123",
		},
		Chat: Chat{
			Id: 3,
			FirstName: "John",
			LastName: "Doe",
			UserName: "johndoe456",
		},
	},
}

func TestGetChatId(t *testing.T) {
	if incomingMessage.GetChatId() != 3 {
		t.Error("Invalid chat ID")
	}
}

func TestGetUser(t *testing.T) {
	user := User{
		Id: 123,
		FirstName: "Jane",
		LastName: "Doe",
		UserName: "janedoe123",
	}

	if incomingMessage.GetUser() != user {
		t.Error("Invalid user")
	}
}

func TestGetMessage(t *testing.T) {
	if incomingMessage.GetMessage() != "Hello world" {
		t.Error("Invalid message")
	}
}
