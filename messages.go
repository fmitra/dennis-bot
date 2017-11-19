package main

import "math/rand"

var messageMap = map[string][]string {
	"basic.track": []string {
		"What are you tracking? You can say something like " +
		"2000JPY for cornerstore sushi. Not in Japan? No problemmm " +
		"you can use any currency!",

		"Let's get started! You can say something like " +
		"4USD for coffee yesterday",
	},
	"basic.help": []string {
		"What?? What do you need help with? Want to do something? Say track!",
		"I can help you track your expenses. Want to start? Just say track!",
	},
	"basic.identity": []string {
		"Dennis Dennis Dennis Dennis",
		"Hiiiiiiiii I'm Dennis!",
	},
	"basic.default": []string {
		"Hi! Tell Dennis what you want to do!",
	},
}

// Returns a message based on a message key. Messages are stored
// as slices for each key and are randomly selected.
func getMessage(messageKey string) (string) {
	messages := messageMap[messageKey]
	totalMessages := len(messages)
	random := rand.Intn(totalMessages)
	return messages[random]
}
