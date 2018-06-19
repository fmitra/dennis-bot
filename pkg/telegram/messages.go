package telegram

// ref: https://core.telegram.org/bots/api#update
// Represents a Telegram Update object. This payload is sent
// to the webhook whenever a user messages us. Message field
// is optional but for our use case this all we care about at
// the moment
type IncomingMessage struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Chat struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

type Message struct {
	MessageId int    `json:"message_id"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
}

// Represents a response to an `IncomingMessage`
type OutgoingMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

// Represents a Telegram Chat Action (ex. Bot is typing indicator)
type ChatAction struct {
	ChatId int `json:"chat_id"`
	Action string `json:"action"`
}

func (i IncomingMessage) GetChatId() (chatId int) {
	return i.Message.Chat.Id
}

func (i IncomingMessage) GetMessage() (message string) {
	return i.Message.Text
}

func (i IncomingMessage) GetUser() User {
	return i.Message.From
}
