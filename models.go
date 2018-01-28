package main

import (
	"errors"
	"log"
	"time"
)

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

// Wit.ai Entity
type WitEntity []struct {
	Value string `json:"value"`
	Confidence float64 `json:"confidence"`
}

// Wit.ai API Response
type WitResponse struct {
	Entities struct {
		Amount WitEntity `json:"amount"`
		DateTime WitEntity `json:"datetime"`
		Description WitEntity `json:"description"`
	} `json:"entities"`
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

func (incMessage IncomingMessage) getChatId() (chatId int) {
	return incMessage.Message.Chat.Id
}

func (incMessage IncomingMessage) getMessage() (message string) {
	return incMessage.Message.Text
}

func (incMessage IncomingMessage) getUser() (User) {
	return incMessage.Message.From
}

// Checks if Wit.ai was able to infer a valid
// Amount Entity from the IncomingMessage.getMessage()
func (witResponse WitResponse) getAmount() (float64, string, error) {
	amount := witResponse.Entities.Amount
	if len(amount) == 0 {
		return 0, "", errors.New("No amount")
	}

	totalAmount, currency := parseAmount(amount[0].Value)

	if totalAmount > 0 && currency != "" {
		return totalAmount, currency, nil
	}

	return 0, "", errors.New("Invalid amount")
}

// Checks if Wit.ai was able to infer a valid
// Description Entity from the IncomingMessage.getMessage()
func (witResponse WitResponse) getDescription() (string, error) {
	description := witResponse.Entities.Description
	if len(description) == 0 {
		return "", errors.New("No description")
	}

	parsedDescription := parseDescription(description[0].Value)
	return parsedDescription, nil
}

// Checks if Wit.ai was able to infer a valid
// Date Entity from the IncomingMessage.getMessage()
// If no date is provided, it will always default to today
func (witResponse WitResponse) getDate() (time.Time) {
	dateTime := witResponse.Entities.DateTime
	stringDate := ""
	if len(dateTime) != 0 {
		stringDate = dateTime[0].Value
	}

	parsedDate := parseDate(stringDate)
	return parsedDate
}

// Infers whether the User is attempting ot track an expense
func (witResponse WitResponse) isTracking() (bool, error) {
	amount, currency, err := witResponse.getAmount()
	if err != nil {
		log.Printf("Cannot infer without amount")
		return false, err
	}

	description, err := witResponse.getDescription()
	if err != nil {
		log.Printf("Cannot infer without description")
		return true, err
	}

	log.Printf("Inferring tracking %s %s %s", amount, currency, description)
	return true, nil
}

// Sets historical USD value of of foriegn expenses
func (expense Expense) setHistorical() {
	if expense.Currency == "USD" {
		expense.Historical = expense.Total
	}
	expense.Historical = float64(0)
}
