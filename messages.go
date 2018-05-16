package main

import "math/rand"

var messageMap = map[string][]string{
	"default": []string{
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
	"track.success": []string{
		"Ok got it!",

		"Roger that!",
	},
	"track.error": []string{
		"I didn't understand that. You need to tell me exactly what " +
			"what your expense is. For example '0.00012BTC for Rent'",
	},
}

// Returns a message based on a message key. Messages are stored
// as slices for each key and are randomly selected.
func getMessage(messageKey string) string {
	messages := messageMap[messageKey]
	totalMessages := len(messages)
	random := rand.Intn(totalMessages)
	return messages[random]
}
