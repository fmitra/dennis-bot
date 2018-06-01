package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/fmitra/dennis/alphapoint"
	"github.com/fmitra/dennis/expenses"
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
	response := bot.BuildResponse(incMessage)

	channel := make(chan bool)

	go func() {
		bot.SendMessage(response, incMessage)
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
func (bot *Bot) SendMessage(response string, incMessage telegram.IncomingMessage) int {
	chatId := incMessage.GetChatId()

	return bot.env.telegram.Send(chatId, response)
}

// Returns a message based on a message key. Messages are stored
// as slices for each key and are randomly selected.
func (bot *Bot) GetMessage(messageKey string, messageVar string) string {
	messages := MessageMap[messageKey]
	totalMessages := len(messages)
	random := rand.Intn(totalMessages)

	var parsedMessage string
	message := messages[random]
	if messageVar != "" && strings.Contains(message, "{{var}}") {
		parsedMessage = strings.Replace(message, "{{var}}", messageVar, -1)
	} else {
		parsedMessage = message
	}

	return parsedMessage
}

func (bot *Bot) BuildResponse(incMessage telegram.IncomingMessage) string {
	w := wit.NewClient(bot.env.config.Wit.Token)
	witResponse := w.ParseMessage(incMessage.GetMessage())
	userId := incMessage.GetUser().Id

	intent := witResponse.GetIntent()

	var err error
	var messageVar string
	var keyword string

	switch intent {
	case wit.TRACKING_SUCCESS:
		go bot.NewExpense(witResponse, userId)
	case wit.PERIOD_TOTAL_SUCCESS:
		messageVar, err = bot.GetTotalByPeriod(witResponse, userId)
	}

	if err != nil {
		keyword = "error"
	} else {
		keyword = intent
	}

	return bot.GetMessage(keyword, messageVar)
}

// Creates an expense item from a Wit.ai response
func (bot *Bot) NewExpense(w wit.WitResponse, userId int) bool {
	date := w.GetDate()
	amount, currency, _ := w.GetAmount()
	description, _ := w.GetDescription()

	a := alphapoint.NewClient(bot.env.config.AlphaPoint.Token)
	historical := a.Convert(
		currency,
		"USD",
		amount,
	)

	expense := &expenses.Expense{
		Date:        date,
		Description: description,
		Total:       amount,
		Historical:  historical,
		Currency:    currency,
		UserId:      userId,
	}
	expenseManager := expenses.NewExpenseManager(bot.env.db)
	return expenseManager.Save(expense)
}

func (bot *Bot) GetTotalByPeriod(response wit.WitResponse, userId int) (string, error) {
	expenseManager := expenses.NewExpenseManager(bot.env.db)
	period, err := response.GetSpendPeriod()
	total, err := expenseManager.TotalByPeriod(period, userId)

	return strconv.FormatFloat(total, 'f', 2, 64), err
}
