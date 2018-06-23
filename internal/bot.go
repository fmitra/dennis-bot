package internal

import (
	"encoding/json"
	"log"

	convo "github.com/fmitra/dennis-bot/internal/conversation"
	"github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/wit"
)

// Bot is responsible for parsing messages and responding
// to a user. It is configured based on the environment
type Bot struct {
	env *Env
}

// Converse is the entry point to communicate with the bot. We parse an incoming
// message and map it to  to a key word trigger to determine a response.
func (bot *Bot) Converse(b []byte) int {
	incMessage, err := bot.ReceiveMessage(b)
	if err != nil {
		log.Printf("bot: cannot respond to unsupported payload - %s", err)
		errorCode := 400
		return errorCode
	}

	bot.SendTypingIndicator(incMessage)
	response := bot.BuildResponse(incMessage)

	return bot.SendMessage(response, incMessage)
}

// ReceiveMessage unmarshals a byte response into a telegram IncomingMessage.
func (bot *Bot) ReceiveMessage(b []byte) (telegram.IncomingMessage, error) {
	var incM telegram.IncomingMessage
	err := json.Unmarshal(b, &incM)
	if err != nil {
		return incM, err
	}

	return incM, nil
}

// SendMessage sends a a message back through Telegram.
func (bot *Bot) SendMessage(r convo.BotResponse, incM telegram.IncomingMessage) int {
	chatID := incM.GetChatID()

	return bot.env.telegram.Send(chatID, string(r))
}

// SendTypingIndicator sends a typign indicator to the user to alert them
// that we have received and processing their mssage.
func (bot *Bot) SendTypingIndicator(incM telegram.IncomingMessage) int {
	chatID := incM.GetChatID()
	action := "typing"
	return bot.env.telegram.SendAction(chatID, action)
}

// BuildResponse coordinates with the action layer to to determine context behind
// a user's message and return an appropriate response.
func (bot *Bot) BuildResponse(incM telegram.IncomingMessage) convo.BotResponse {
	w := wit.NewClient(bot.env.config.Wit.Token)
	witResponse := w.ParseMessage(incM.GetMessage())
	actions := &convo.Actions{
		Db:     bot.env.db,
		Cache:  bot.env.cache,
		Config: bot.env.config,
	}
	botResponse := convo.GetResponse(witResponse, incM, actions)
	return botResponse
}
