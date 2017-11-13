package main

import (
	"log"
	"encoding/json"
	"strings"
)

func converse(payload []byte) {
	incMessage, err := receiveMessage(payload)
	if err != nil {
		panic(err)
	}
	keyword := mapToKeyword(incMessage)
	sendMessage(keyword, incMessage.getChatId())
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
func sendMessage(keyword string, chatId int) {
	var message string
	switch keyword {
	case "track":
		message = Messages["basic.track"]
	case "help":
		message = Messages["basic.help"]
	case "identity":
		message = Messages["basic.identity"]
	default:
		message = Messages["basic.default"]
	}

	go telegram.send(chatId, message)
}

// IncomingMessages are mapped to keywords to trigger the approriate
// message for a user's intent.
func mapToKeyword(incMessage IncomingMessage) (string) {
	var responseTriggers = map[string][]string {
		"default": []string{
			"hello",
		},
		"help": []string {
			"help",
		},
		"identity": []string {
			"who are",
			"your name",
		},
		"track": []string{
			"expense",
			"track",
			"log",
		},
	}

	message := incMessage.getMessage()
	for keyword, triggers := range responseTriggers {
		for _, trigger := range triggers {
			parsedMessage := strings.ToLower(message)
			if strings.Contains(parsedMessage, trigger) {
				return keyword
			}
		}
	}
	return "default"
}
