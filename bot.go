package main

import (
	"log"
	"encoding/json"

	"github.com/fmitra/dennis/wit"
	"github.com/fmitra/dennis/expenses"
	"github.com/fmitra/dennis/telegram"
	"github.com/fmitra/dennis/users"
)

// Entry point to communicate with Dennis.
// We parse incoming messages from telegram and
// map it to key word triggers to determine a response
// to send back to the user.
func converse(payload []byte) {
	incMessage, err := receiveMessage(payload)
	if err != nil {
		panic(err)
	}
	users.UpdateSession(incMessage.GetUser())
	keyword := mapToKeyword(incMessage)
	sendMessage(keyword, incMessage)
}

// Parses an incoming message to ensure valid JSON
// was returned.
func receiveMessage(payload []byte) (telegram.IncomingMessage, error) {
	var incMessage telegram.IncomingMessage
	err := json.Unmarshal(payload, &incMessage)
	if err != nil {
		log.Printf("bot: invalid payload - %s", err)
		return telegram.IncomingMessage{}, err
	}

	return incMessage, nil
}

// Determines outgoing response based on the keyword
func sendMessage(keyword string, incMessage telegram.IncomingMessage) {
	message := getMessage(keyword)
	chatId := incMessage.GetChatId()
	go telegram.Client.Send(chatId, message)
}

// IncomingMessages are mapped to keywords to trigger the approriate
// message for a user's intent.
func mapToKeyword(incMessage telegram.IncomingMessage) (string) {
	witResponse := wit.Client.ParseMessage(incMessage.GetMessage())
	isTracking, err := witResponse.IsTracking()
	if isTracking == true && err == nil {
		log.Printf("%s", witResponse)
		go expenses.NewExpense(witResponse)
		return "track.success"
	}

	if isTracking == true && err != nil {
		return "track.error"
	}

	return "default"
}
