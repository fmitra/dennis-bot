package main

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
		From struct {
			Id int `json:"id"`
			FirstName string `json:"first_name"`
			LastName string `json:"last_name"`
			UserName string `json:"username"`
		} `json:"from"`
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

// Describes a tracked expense
type Expense struct {
	Date string         // Date the expense was made
	Description string  // Description of the expense
	Total float64       // Total amount paid for the expense
	Historical float64  // Historical USD value of the total
	Currency string     // Currency denomination of the total
	Category string     // Category of the expense
}

func (incMessage IncomingMessage) getChatId() (int) {
	return incMessage.Message.Chat.Id
}

func (incMessage IncomingMessage) getMessage() (string) {
	return incMessage.Message.Text
}

// Sets historical USD value of of foriegn expenses
func (expense Expense) setHistorical() {
	if expense.Currency == "USD" {
		expense.Historical = expense.Total
	}
	expense.Historical = float64(0)
}

// Parses a line of text to build an Expense.
//
// For example:
//   >> "6USD for Locavore 12/15/2017"
//   >> "50SGD for DTF yesterday"
//   >> "20000JPY for Sushi"
func parseExpense(expense string) (Expense) {
	return Expense{}
}
