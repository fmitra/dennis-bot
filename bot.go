package main

import (
	"encoding/json"
	"log"
	"strconv"
	"math/rand"

	"github.com/fmitra/dennis/alphapoint"
	"github.com/fmitra/dennis/sessions"
	"github.com/fmitra/dennis/telegram"
	"github.com/fmitra/dennis/wit"
)

// Bot is responsible for parsing messages and responding
// to a user. It is configured based on the environment
type Bot struct {
	env *Env
}

// Entry point to communicate with the bot. We parse an incoming message
// and map it to  to a key word trigger to determine a response
func (bot *Bot) Converse(payload []byte) {
	incMessage, err := bot.ReceiveMessage(payload)
	if err != nil {
		panic(err)
	}
	user := incMessage.GetUser()
	sessions.Set(strconv.Itoa(user.Id), user)
	keyword := bot.MapToKeyword(incMessage)
	bot.SendMessage(keyword, incMessage)
}

// Parses an incoming telegram message
func (bot *Bot) ReceiveMessage(payload []byte) (telegram.IncomingMessage, error) {
	var incMessage telegram.IncomingMessage
	err := json.Unmarshal(payload, &incMessage)
	if err != nil {
		log.Printf("bot: invalid payload - %s", err)
		return telegram.IncomingMessage{}, err
	}

	return incMessage, nil
}

// Sends a message back through telegram
func (bot *Bot) SendMessage(keyword string, incMessage telegram.IncomingMessage) {
	message := bot.GetResponse(keyword)
	chatId := incMessage.GetChatId()
	go telegram.Client.Send(chatId, message)
}

// Returns a message based on a message key. Messages are stored
// as slices for each key and are randomly selected.
func (bot *Bot) GetResponse(messageKey string) string {
	messages := messageMap[messageKey]
	totalMessages := len(messages)
	random := rand.Intn(totalMessages)
	return messages[random]
}

// IncomingMessages are mapped to keywords to trigger the approriate
// message for a user's intent.
func (bot *Bot) MapToKeyword(incMessage telegram.IncomingMessage) string {
	witResponse := wit.Client.ParseMessage(incMessage.GetMessage())
	isTracking, err := witResponse.IsTracking()
	if isTracking == true && err == nil {
		log.Printf("%s", witResponse)
		go bot.NewExpense(witResponse, incMessage.GetUser().Id)
		return "track.success"
	}

	if isTracking == true && err != nil {
		return "track.error"
	}

	return "default"
}

// Creates an expense item from a Wit.ai response
func (bot *Bot) NewExpense(w wit.WitResponse, userId int) {
	date := w.GetDate()
	amount, currency, _ := w.GetAmount()
	description, _ := w.GetDescription()
	historical := alphapoint.Client.Convert(
		currency,
		"USD",
		amount,
	)

	bot.env.db.Create(&Expense{
		Date:        date,
		Description: description,
		Total:       amount,
		Historical:  historical,
		Currency:    currency,
		UserId:      userId,
	})
}
