package internal

import (
	"log"
	"math/rand"
	"strings"
)

const (
	TRACKING_SUCCESS = "tracking_success"

	TRACKING_ERROR = "tracking_error"

	PERIOD_TOTAL_SUCCESS = "period_total_success"

	DEFAULT = "default"

	ERROR = "error"
)

type BotResponse string

var MessageMap = map[string][]string{
	DEFAULT: []string{
		"Hi! Tell Dennis what you want to do!",

		"What are you tracking? You can say something like " +
			"2000JPY for cornerstore sushi. Not in Japan? No problemmm " +
			"you can use any currency!",

		"Let's get started! You can say something like " +
			"4USD for coffee yesterday",

		"Dennis Dennis Dennis Dennis",

		"Hiiiiiiiii I'm Dennis!",

		"How much did you spend? You can say something like 20USD " +
			"for Dinner",
	},
	ERROR: []string{
		"Sorry can't answer. I'm going out for lunch",

		"Ask me something later, Dennis is on vacation",
	},
	TRACKING_SUCCESS: []string{
		"Ok got it!",

		"Roger that!",
	},
	TRACKING_ERROR: []string{
		"I didn't understand that. You need to tell me exactly what " +
			"what your expense is. For example '0.00012BTC for Rent'",

		"I didn't get that. Try saying 'How much did I spend this week?'",
	},
	PERIOD_TOTAL_SUCCESS: []string{
		"You spent {{var}}",
	},
}

// Returns a message based on a message key. Messages are stored
// as slices for each key and are randomly selected.
func GetMessage(messageKey string, messageVar string) BotResponse {
	messages := MessageMap[messageKey]
	totalMessages := len(messages)
	log.Printf("%s hmm %s %s", totalMessages, messageKey, messages)
	random := rand.Intn(totalMessages)

	var parsedMessage string
	message := messages[random]
	if messageVar != "" && strings.Contains(message, "{{var}}") {
		parsedMessage = strings.Replace(message, "{{var}}", messageVar, -1)
	} else {
		parsedMessage = message
	}

	return BotResponse(parsedMessage)
}
