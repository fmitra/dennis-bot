package conversation

import (
	"math/rand"
	"strings"
)

const (
	// Generic default responses when the bot is unable to determine a
	// Users intent
	DEFAULT = "default"

	// Requests a password for the user to start interacting with the bot
	ONBOARD_USER_ASK_FOR_PASSWORD = "onboard_user_ask_for_password"

	// Response when a user rejects password confirmation
	ONBOARD_USER_REJECT_PASSWORD = "onboard_user_reject_password"

	// Response when a user submits a password. We must first get a second
	// confirmation before saving it.
	ONBOARD_USER_CONFIRM_PASSWORD = "onboard_user_confirm_password"

	// Response when a user responds to a password confirmation prompt with
	// invalid text
	ONBOARD_USER_CONFIRM_PASSWORD_ERROR = "onboard_user_confirm_password_error"

	// Response when a user successfully sets a password
	ONBOARD_USER_SAY_OUTRO = "onboard_user_say_outro"

	// Response when the bot successfully tracks an expense
	TRACK_EXPENSE_SUCCESS = "track_expense_success"

	// Response when a tracking expense request failed
	TRACK_EXPENSE_ERROR = "track_expense_error"

	// Respone when a user requests for expense total by time period
	GET_EXPENSE_TOTAL_SUCCESS = "get_expense_total_success"

	// Response when an expense total query failed
	GET_EXPENSE_TOTAL_ERROR = "get_expense_total_error"
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

		"How much did you spend? You can say something like 20USD for Dinner",
	},
	GET_EXPENSE_TOTAL_ERROR: []string{
		"I didn't understand that. You can say something like " +
			"'how much did I spend today?'",
	},
	GET_EXPENSE_TOTAL_SUCCESS: []string{
		"You spent {{var}}",
	},
	TRACK_EXPENSE_SUCCESS: []string{
		"Ok got it!",

		"Roger that!",
	},
	TRACK_EXPENSE_ERROR: []string{
		"I didn't understand that. You need to tell me exactly what " +
			"what your expense is. For example '0.00012BTC for Rent'",

		"I didn't get that. Try saying 'How much did I spend this week?'",
	},
	ONBOARD_USER_ASK_FOR_PASSWORD: []string{
		"You seem around here! My name is Dennis and I can track your expenses " +
			"but first, you need to create a password to protect your data. What " +
			"do you want your password to be?",

		"You need to create a password before I can track your expenses. What " +
			"do you want your password to be?",
	},
	ONBOARD_USER_CONFIRM_PASSWORD: []string{
		"Ok your password is {{var}}. Is that right? Just say yes or no.",
	},
	ONBOARD_USER_SAY_OUTRO: []string{
		"Alright you're all set! When you're ready to start, you can say something " +
			"like 450RUB for Lunch",
	},
	ONBOARD_USER_CONFIRM_PASSWORD_ERROR: []string{
		"I didn't understand that. You can say yes or no",
	},
	ONBOARD_USER_REJECT_PASSWORD: []string{
		"Okay will if you think of password later just let me know.",
	},
}

// Returns a message based on a message key. Messages are stored
// as slices for each key and are randomly selected.
func GetMessage(messageKey string, messageVar string) BotResponse {
	messages := MessageMap[messageKey]
	totalMessages := len(messages)
	random := rand.Intn(totalMessages)

	message := messages[random]
	parsedMessage := message
	if messageVar != "" && strings.Contains(message, "{{var}}") {
		parsedMessage = strings.Replace(message, "{{var}}", messageVar, -1)
	}

	return BotResponse(parsedMessage)
}
