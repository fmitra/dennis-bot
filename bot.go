package main

import (
	"log"
	"encoding/json"
)

func converse(payload []byte) {
	incMessage, err := receiveMessage(payload)
	if err != nil {
		panic(err)
	}
	updateSession(incMessage.getUser())
	keyword := mapToKeyword(incMessage)
	sendMessage(keyword, incMessage)
}

// Parses an incoming message to ensure valid JSON
// was returned.
func receiveMessage(payload []byte) (IncomingMessage, error) {
	var incMessage IncomingMessage
	err := json.Unmarshal(payload, &incMessage)
	if err != nil {
		log.Printf("bot: invalid payload - %s", err)
		return IncomingMessage{}, err
	}

	return incMessage, nil
}

// Determines outgoing response based on the keyword
func sendMessage(keyword string, incMessage IncomingMessage) {
	message := getMessage(keyword)
	chatId := incMessage.getChatId()
	go telegram.send(chatId, message)
}

// IncomingMessages are mapped to keywords to trigger the approriate
// message for a user's intent.
func mapToKeyword(incMessage IncomingMessage) (string) {
	witResponse := witAi.parseMessage(incMessage.getMessage())
	isTracking, err := witResponse.isTracking()
	if isTracking == true && err == nil {
		return "track.success"
	}

	if isTracking == true && err != nil {
		return "track.error"
	}

	return "default"
}
