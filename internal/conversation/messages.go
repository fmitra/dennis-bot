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

	// Response when we are unable to create an account during the onboarding
	// flow
	ONBOARD_USER_ACCOUNT_CREATION_FAILED = "onboard_user_account_creation_failed"

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

	// Response when a user first requests expense totals
	GET_EXPENSE_TOTAL_ASK_FOR_PASSWORD = "get_expense_total_ask_for_password"

	// Response when a user provides an invalid password
	GET_EXPENSE_TOTAL_PASSWORD_INVALID = "get_expense_total_password_invalid"

	// Response when an invalid time range is requested by the user
	GET_EXPENSE_TOTAL_INVALID_PERIOD = "get_expense_total_invalid_period"
)

type BotResponse string

var MessageMap = map[string][]string{
	DEFAULT: []string{
		"alright, alright you want my attention? You have it",

		"I just got back from accounting school. Tell me to track something. " +
			"You can say '2000JPY for sushi'",

		"hmm... ehhhh.... idk",

		"H-h-hi I'm dennis",

		"What you want!? I'm trying to take a vacation",
	},
	GET_EXPENSE_TOTAL_ERROR: []string{
		"ehhh oh no... Idk what happened. Try again later.",
	},
	GET_EXPENSE_TOTAL_SUCCESS: []string{
		"You spent {{var}} USD",

		"lets see now... looks like {{var}} USD",

		"hmmm.... {{var}} USD",
	},
	GET_EXPENSE_TOTAL_INVALID_PERIOD: []string{
		"I'm not that smart. Ask me something like 'how much did I spend today'",

		"ask me something later. Dennis is on vacation",
	},
	GET_EXPENSE_TOTAL_ASK_FOR_PASSWORD: []string{
		"wait wait! What's your password",

		"I can't tell you without a password",

		"ok, but first you gotta tell me your password",
	},
	GET_EXPENSE_TOTAL_PASSWORD_INVALID: []string{
		"wat! that password's not rigiht. Try again or say 'cancel'",

		"hold up! that's not the password. Try again or say 'cancel'",
	},
	TRACK_EXPENSE_SUCCESS: []string{
		"ok writing it down...",

		"you spend so much. When you gon take me out? today is dennis day!",

		"okay one min. Let me get my calculator",

		"writing.... and... done!",
	},
	TRACK_EXPENSE_ERROR: []string{
		"I didn't get that. Try saying something like '1000RUB for lunch'",

		"hmm this is embarassing. I have no idea what I'm doing. Try asking me to track " +
		"12USD for food",
	},
	ONBOARD_USER_ASK_FOR_PASSWORD: []string{
		"I dont know you. I'm Dennis though, and I can track your finances. But first, " +
		"you gotta make a password. What do you want your password to be?",

		"Hiiiii I'm Dennis. I track your expenses. What do you want your password to be?",

		"h-h-hi I'm dennis. I can track your finances, but first, we got to make a " +
		"password. What do you want your password to be?",
	},
	ONBOARD_USER_CONFIRM_PASSWORD: []string{
		"alright got it, your password is '{{var}}'. Is that right? Just say yes or no.",

		"ok '{{var}}' right? Just say yes or no",
	},
	ONBOARD_USER_SAY_OUTRO: []string{
		"got it! you're all set! Next time you buy something, just to tell me something " +
		"like 450SGD for tickets",
	},
	ONBOARD_USER_CONFIRM_PASSWORD_ERROR: []string{
		"I don't get it. I told you just say yes or no!",

		"You're confusing me. Just say 'yes' or 'no'",
	},
	ONBOARD_USER_REJECT_PASSWORD: []string{
		"alright well just come back later if you think of a password",
	},
	ONBOARD_USER_ACCOUNT_CREATION_FAILED: []string{
		"Hey! I think you already have an account. Why don't you try tracking something",
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
