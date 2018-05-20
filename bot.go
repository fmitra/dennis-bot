package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"

	"github.com/fmitra/dennis/alphapoint"
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
func (bot *Bot) Converse(payload []byte) chan bool {
	incMessage, err := bot.ReceiveMessage(payload)
	if err != nil {
		panic(err)
	}
	user := incMessage.GetUser()
	bot.env.cache.Set(strconv.Itoa(user.Id), user)
	keyword := bot.MapToKeyword(incMessage)

	channel := make(chan bool)

	go func() {
		bot.SendMessage(keyword, incMessage)
		channel <- true
		close(channel)
	}()

	return channel
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
func (bot *Bot) SendMessage(keyword string, incMessage telegram.IncomingMessage) int {
	message := bot.GetResponse(keyword)
	chatId := incMessage.GetChatId()

	return bot.env.telegram.Send(chatId, message)
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
	w := wit.Client(bot.env.config.Wit.Token)
	witResponse := w.ParseMessage(incMessage.GetMessage())
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
func (bot *Bot) NewExpense(w wit.WitResponse, userId int) bool {
	date := w.GetDate()
	amount, currency, _ := w.GetAmount()
	description, _ := w.GetDescription()

	a := alphapoint.Client(bot.env.config.AlphaPoint.Token)
	historical := a.Convert(
		currency,
		"USD",
		amount,
	)

	expense := &Expense{
		Date:        date,
		Description: description,
		Total:       amount,
		Historical:  historical,
		Currency:    currency,
		UserId:      userId,
	}
	expenseManager := NewExpenseManager(bot.env.db)
	return expenseManager.Save(expense)
}
