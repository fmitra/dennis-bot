package telegram

// User is a user with a Telegram chat account who initiates contact
// with the bot service. While any Telegram user may interact with the bot,
// actual features are limited to users who create accounts with the bot.
type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}
