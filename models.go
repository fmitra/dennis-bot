package main

type User struct {
	Id int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	UserName string `json:"username"`
}

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
		From User `json:"from"`
		Chat struct {
			Id int `json:"id"`
			FirstName string `json:"first_name"`
			LastName string `json:"last_name"`
			UserName string `json:"username"`
		} `json:"chat"`
	} `json:"message"`
}

// Bot response to an `IncomingMessage`
type OutgoingMessage struct {
	ChatId int `json:"chat_id"`
	Text string `json:"text"`
}

func (incMessage IncomingMessage) getChatId() (chatId int) {
	return incMessage.Message.Chat.Id
}

func (incMessage IncomingMessage) getMessage() (message string) {
	return incMessage.Message.Text
}

func (incMessage IncomingMessage) getUser() (User) {
	return incMessage.Message.From
}
