package telegram

// Represents a User interacting with the bot
type User struct {
	Id int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	UserName string `json:"username"`
}
