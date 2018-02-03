package telegram

import (
	"github.com/fmitra/dennis/users"
)

// ref: https://core.telegram.org/bots/api#update
// Represents a Telegram Update object. This payload is sent
// to the webhook whenever a user messages us. Message field
// is optional but for our use case this all we care about at
// the moment
type IncomingMessage struct {
	UpdateId int `json:"update_id"`
	Message struct {
		MessageId int `json:"message_id"`
		Date int `json:"date"`
		Text string `json:"text"`
		From users.User `json:"from"`
		Chat struct {
			Id int `json:"id"`
			FirstName string `json:"first_name"`
			LastName string `json:"last_name"`
			UserName string `json:"username"`
		} `json:"chat"`
	} `json:"message"`
}

// Represents a response to an `IncomingMessage`
type OutgoingMessage struct {
	ChatId int `json:"chat_id"`
	Text string `json:"text"`
}

func (i IncomingMessage) GetChatId() (chatId int) {
	return i.Message.Chat.Id
}

func (i IncomingMessage) GetMessage() (message string) {
	return i.Message.Text
}

func (i IncomingMessage) GetUser() (users.User) {
	return i.Message.From
}
