package conversation

import (
	"math/rand"
	"strings"
)

const (
	// DefaultResponse is a generic response when the bot is unable to determine
	// a Users intent
	DefaultResponse = "default"

	// OnboardUserAskForPassword requests a password for the user to start
	// interacting with the bot
	OnboardUserAskForPassword = "onboard_user_ask_for_password"

	// OnboardUserRejectPassword is response when a user rejects password confirmation
	OnboardUserRejectPassword = "onboard_user_reject_password"

	// OnboardUserConfirmPassword is a response when a user submits a password. We must
	// first get a second confirmation before saving it.
	OnboardUserConfirmPassword = "onboard_user_confirm_password"

	// OnboardUserConfirmPasswordError is a response when a user responds to a password
	// confirmation prompt with invalid text.
	OnboardUserConfirmPasswordError = "onboard_user_confirm_password_error"

	// OnboardUserConfirmPasswordInvalid is a response when we receive an error while
	// encrypting a password.
	OnboardUserConfirmPasswordInvalid = "onboard_user_confirm_password_invalid"

	// OnboardUserAccountCreationFailed is a rsponse when we are unable to create
	// an account during the onboarding flow
	OnboardUserAccountCreationFailed = "onboard_user_account_creation_failed"

	// OnboardUserDecryptionFailed is a response when we fail to decrypt the user password
	// from cache.
	OnboardUserDecryptionFailed = "onboard_user_decryption_failed"

	// OnboardUserAskForCurrency is a response to check what currency the user would
	// like to receive expense history in.
	OnboardUserAskForCurrency = "onboard_user_ask_for_currency"

	// OnboardUserInvalidCurrency is a response when a user does not it give us a valid
	// currency ISO.
	OnboardUserInvalidCurrency = "onboard_user_invalid_currency"

	// OnboardUserSayOutro is sent when a user successfully sets a password
	OnboardUserSayOutro = "onboard_user_say_outro"

	// TrackExpenseSuccess is a response when the bot successfully tracks an expense
	TrackExpenseSuccess = "track_expense_success"

	// TrackExpenseError is a response when a tracking expense request failed
	TrackExpenseError = "track_expense_error"

	// GetExpenseTotalSuccess is a response when a user requests for expense total
	// by time period
	GetExpenseTotalSuccess = "get_expense_total_success"

	// GetExpenseTotalError is a response when an expense total query failed
	GetExpenseTotalError = "get_expense_total_error"

	// GetExpenseTotalAskForPassword is a response when a user first requests expense totals
	GetExpenseTotalAskForPassword = "get_expense_total_ask_for_password"

	// GetExpenseTotalPasswordInvalid is a response when a user provides an invalid password
	GetExpenseTotalPasswordInvalid = "get_expense_total_password_invalid"

	// GetExpenseTotalInvalidPeriod is a response when an invalid time range is requested
	// by the user
	GetExpenseTotalInvalidPeriod = "get_expense_total_invalid_period"

	// GetExpenseTotalCancel is a response when the user requests the bot to cancel
	// the query
	GetExpenseTotalCancel = "get_expense_total_cancel"
)

// BotResponse is a message delivered to the User from the Bot
type BotResponse string

// MessageMap contains all messages to that may be sent to a user based
// on a key word.
var MessageMap = map[string][]string{
	DefaultResponse: []string{
		"alright, alright you want my attention? You have it",

		"I just got back from accounting school. Tell me to track something. " +
			"You can say '2000JPY for sushi'",

		"hmm... ehhhh.... idk",

		"H-h-hi I'm dennis",

		"What you want!? I'm trying to take a vacation",
	},
	GetExpenseTotalError: []string{
		"ehhh oh no... Idk what happened. Try again later.",
	},
	GetExpenseTotalSuccess: []string{
		"You spent {{var}}",

		"lets see now... looks like {{var}}",

		"hmmm.... {{var}}",
	},
	GetExpenseTotalInvalidPeriod: []string{
		"I'm not that smart. Ask me something like 'how much did I spend today'",

		"ask me something later. Dennis is on vacation",
	},
	GetExpenseTotalAskForPassword: []string{
		"wait wait! What's your password",

		"I can't tell you without a password",

		"ok, but first you gotta tell me your password",
	},
	GetExpenseTotalPasswordInvalid: []string{
		"wat! that password's not rigiht. Try again or say 'cancel'",

		"hold up! that's not the password. Try again or say 'cancel'",
	},
	GetExpenseTotalCancel: []string{
		"ok, just message me if you change your mind later.",
	},
	TrackExpenseSuccess: []string{
		"ok writing it down...",

		"you spend so much. When you gon take me out? today is dennis day!",

		"okay one min. Let me get my calculator",

		"writing.... and... done!",
	},
	TrackExpenseError: []string{
		"I didn't get that. Try saying something like '1000RUB for lunch'",

		"hmm this is embarrassing. I have no idea what I'm doing. Try asking me to track " +
			"12USD for food",
	},
	OnboardUserAskForPassword: []string{
		"I dont know you. I'm Dennis though, and I can track your finances. But first, " +
			"you gotta make a password. What do you want your password to be?",

		"Hiiiii I'm Dennis. I track your expenses. What do you want your password to be?",

		"h-h-hi I'm dennis. I can track your finances, but first, we got to make a " +
			"password. What do you want your password to be?",
	},
	OnboardUserConfirmPassword: []string{
		"alright got it, your password is '{{var}}'. Is that right? Just say yes or no.",

		"ok '{{var}}' right? Just say yes or no",
	},
	OnboardUserAskForCurrency: []string{
		"what currency do you want to receive updates in? You can say something like 'USD' " +
			"or 'SGD' or any other currency ISO",
	},
	OnboardUserInvalidCurrency: []string{
		"hey I don't understand that! Please say a currency ISO like 'USD' or 'JPY'",
	},
	OnboardUserSayOutro: []string{
		"got it! you're all set! Next time you buy something, just to tell me something " +
			"like 450SGD for tickets",
	},
	OnboardUserDecryptionFailed: []string{
		"whoops something went wrong with your password. Let's start over.",
	},
	OnboardUserConfirmPasswordError: []string{
		"I don't get it. I told you just say yes or no!",

		"You're confusing me. Just say 'yes' or 'no'",
	},
	OnboardUserConfirmPasswordInvalid: []string{
		"I can't use this as a password, try something else",
	},
	OnboardUserRejectPassword: []string{
		"alright well just come back later if you think of a password",
	},
	OnboardUserAccountCreationFailed: []string{
		"Hey! I think you already have an account. Why don't you try tracking something",
	},
}

// GetMessage returns a message based on a message key. Messages are stored
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
