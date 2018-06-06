package internal

import (
	"encoding/json"
	"log"
	"strconv"

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
func (bot *Bot) Converse(payload []byte) int {
	incMessage, err := bot.ReceiveMessage(payload)
	if err != nil {
		log.Panicf("bot: cannot respond to unsupported payload - %s", err)
	}
	user := incMessage.GetUser()
	bot.env.cache.Set(strconv.Itoa(user.Id), user)
	response := bot.BuildResponse(incMessage)

	return bot.SendMessage(response, incMessage)
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
func (bot *Bot) SendMessage(response BotResponse, incMessage telegram.IncomingMessage) int {
	chatId := incMessage.GetChatId()

	return bot.env.telegram.Send(chatId, string(response))
}

// Processes an incoming message and retrieves the appropriate response
func (bot *Bot) BuildResponse(incMessage telegram.IncomingMessage) BotResponse {
	w := wit.NewClient(bot.env.config.Wit.Token)
	witResponse := w.ParseMessage(incMessage.GetMessage())
	userId := incMessage.GetUser().Id

	actions := &Actions{
		env:         bot.env,
		witResponse: witResponse,
		userId:      userId,
	}

	// Incoming messages are associated with user intent. Here we
	// we process the user's request and retrieve a response
	botResponse := actions.ProcessIntent()

	return botResponse
}
