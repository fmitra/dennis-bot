package telegram

// IncomingMessage is a payload sent from Telegram representing
// a User's message.
type IncomingMessage struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Chat contains User data related to a Message.
type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

// Message contains the the core content behind an
// IncomingMessage.
type Message struct {
	MessageID int    `json:"message_id"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
}

// OutgoingMessage is a response to an IncomingMessage.
type OutgoingMessage struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

// ChatAction is a UI indicator that the bot is preparing content, for
// example a typing indicator or loading indicator.
type ChatAction struct {
	ChatID int    `json:"chat_id"`
	Action string `json:"action"`
}

// GetChatID returns the ID of an IncomingMessage.
func (i IncomingMessage) GetChatID() (chatID int) {
	return i.Message.Chat.ID
}

// GetMessage returns the text content of an IncomingMessaeg.
func (i IncomingMessage) GetMessage() (message string) {
	return i.Message.Text
}

// GetUser returns the sender of an IncomingMessage.
func (i IncomingMessage) GetUser() User {
	return i.Message.From
}
