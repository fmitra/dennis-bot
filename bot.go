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
	var message string
	switch keyword {
	case "dogfood":
		message = getMessage("dogfood")
	case "track":
		message = getIntentResponse(keyword, incMessage)
	case "help":
		message = getMessage("basic.help")
	case "identity":
		message = getMessage("basic.identity")
	default:
		message = getMessage("basic.default")
	}

	chatId := incMessage.getChatId()
	go telegram.send(chatId, message)
}

// IncomingMessages are mapped to keywords to trigger the approriate
// message for a user's intent.
func mapToKeyword(incMessage IncomingMessage) (string) {
	witAi.parseMessage(incMessage.getMessage())

	// If a session is already in progress, return the keyword of the
	// session to allo the conversation to continue
	intentSession, err := getIntentSession(incMessage.getUser().Id)
	if err == nil {
		return intentSession.Keyword
	}

	var responseTriggers = map[string][]string {
		"dogfood": []string {
			"dogfood",
		},
		"default": []string {
			"hello",
		},
		"help": []string {
			"help",
		},
		"identity": []string {
			"who are",
			"your name",
		},
		"track": []string {
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
