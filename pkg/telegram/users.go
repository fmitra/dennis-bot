package telegram

// Represents a Telegram user interacting with the bot
// While any Telegram user may interact with the bot, they can
// only use features after setting up an account which is linked
// to a password and a unique ID, separate from their Telegram user ID
type User struct {
	Id        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}
