package conversation

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/fmitra/dennis-bot/pkg/sessions"
	t "github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/users"
	"github.com/fmitra/dennis-bot/pkg/wit"
)

const (
	// Onboards a user into the bot by creating a user account
	ONBOARD_USER_INTENT = "onboard_user_intent"

	// Tracks a user's expense
	TRACK_EXPENSE_INTENT = "track_expense_intent"

	// Returns total expense history for the user
	GET_EXPENSE_TOTAL_INTENT = "get_expense_total_intent"
)

// Intents describe what a user wants to accomplish from a conversation.
// They are responsible for the business logic behind any action we
// take in response to a user's incoming message and return a response
// for that action.
type Intent interface {
	// Determine's a user's intent behind an incoming message
	Respond() (BotResponse, *Context)

	// Returns a list of functions to deliver responses to a user.
	// Functions are listed in order of delivery from top to bottom.
	GetResponses() []func() (BotResponse, error)
}

// Context is embeded into all Intents to provides necessary info and
// methods to process a response
type Context struct {
	Step        int
	WitResponse wit.WitResponse
	IncMessage  t.IncomingMessage
	BotUserId   uint
	AuxData     string
}

// Represents a user's place in a conversation, for example a conversation
// may have just recently been initialized or may be ongoing.
type Conversation struct {
	Intent    string // The objective the of the user. Conversations are based around intents
	UserId    uint   // The user ID from the incoming chat service (ex. Telegram User ID)
	BotUserId uint   // The ID of the user account (if it exists) in our system
	Step      int    // A user's place in a conversation
	AuxData   string // Optional auxiliary info that we can set while processing a response
}

// Helper method to ensure all Intents embeding context has access to the same
// response logic. Intent response methods should return a BotResponse and an error.
// Nil errors and empty responses are typically performed in validation steps (ex.
// ask a question and check for an answer). Valid responses increment a step,
// allowing a user to receive the next response in defined in the `GetResponses` array
func (cx *Context) Process(responses []func() (BotResponse, error)) (BotResponse, *Context) {
	response, err := responses[cx.Step]()

	for response == BotResponse("") && err == nil {
		cx.Step += 1
		response, err = responses[cx.Step]()
	}

	if err == nil && cx.Step != -1 {
		cx.Step += 1
	}

	if cx.Step >= len(responses) {
		cx.EndConversation()
	}

	return response, cx
}

// Skips over to the next response in line
func (cx *Context) SkipResponse() (BotResponse, error) {
	return BotResponse(""), nil
}

// Ends a conversation
func (cx *Context) EndConversation() {
	finalStep := -1
	cx.Step = finalStep
}

func (c *Conversation) HasResponse() bool {
	finalStep := -1
	return c.Step > finalStep
}

// Creates a new intent with additional context fields in order to
// determine a response
func (c *Conversation) GetIntent(w wit.WitResponse, inc t.IncomingMessage, a *Actions, uid uint) Intent {
	context := Context{
		WitResponse: w,
		IncMessage:  inc,
		Step:        c.Step,
		BotUserId:   uid,
		AuxData:     c.AuxData,
	}
	switch c.Intent {
	case ONBOARD_USER_INTENT:
		return &OnboardUser{context, a}
	case TRACK_EXPENSE_INTENT:
		return &TrackExpense{context, a}
	case GET_EXPENSE_TOTAL_INTENT:
		return &GetExpenseTotal{context, a}
	default:
		return &GenericResponse{context, a}
	}
}

// Returns a response to the user. Responses are controlled by Intent types
// which we create on demand with the necessary Context around the ongoing dialog
func (c *Conversation) Respond(w wit.WitResponse, inc t.IncomingMessage, a *Actions) BotResponse {
	if !c.HasResponse() {
		return BotResponse("")
	}

	intent := c.GetIntent(w, inc, a, c.BotUserId)
	response, context := intent.Respond()
	c.Step = context.Step
	c.AuxData = context.AuxData

	return response
}

func InferIntent(w wit.WitResponse, botUserId uint) string {
	noId := 0
	// We can only process a user's request if they have an account
	// in our system, otherwise we force them into an onboarding flow.
	if int(botUserId) == noId {
		return ONBOARD_USER_INTENT
	}

	overview := w.GetMessageOverview()
	switch overview {
	case wit.TRACKING_REQUESTED_SUCCESS:
		return TRACK_EXPENSE_INTENT
	case wit.TRACKING_REQUESTED_ERROR:
		return TRACK_EXPENSE_INTENT
	case wit.EXPENSE_TOTAL_REQUESTED_SUCCESS:
		return GET_EXPENSE_TOTAL_INTENT
	default:
		return ""
	}
}

// Creates a new conversation between the bot and the user. If the user
// has an existing account, we associate their account ID with the user
// ID of their chat service.
func NewConversation(userId uint, w wit.WitResponse, a *Actions) Conversation {
	var botUser users.User
	a.Db.Where("telegram_id = ?", userId).First(&botUser)

	intent := InferIntent(w, botUser.ID)
	conversation := Conversation{
		Intent:    intent,
		UserId:    userId,
		BotUserId: botUser.ID,
	}

	return conversation
}

func GetConversation(userId uint, cache sessions.Session) (Conversation, error) {
	var conversation Conversation
	cacheKey := fmt.Sprintf("%s_conversation", strconv.Itoa(int(userId)))
	cache.Get(cacheKey, &conversation)
	noUser := 0
	if int(conversation.UserId) == noUser {
		return conversation, errors.New("No conversation found")
	}
	if !conversation.HasResponse() {
		return conversation, errors.New("No responses available")
	}
	return conversation, nil
}

func GetResponse(w wit.WitResponse, inc t.IncomingMessage, a *Actions) BotResponse {
	userId := inc.GetUser().Id

	conversation, err := GetConversation(userId, a.Cache)
	if err != nil {
		conversation = NewConversation(userId, w, a)
	}

	response := conversation.Respond(w, inc, a)

	// Check if there are additional responses available. If responses are found,
	// we cache the conversation so the user can pick up where they left off
	cacheKey := fmt.Sprintf("%s_conversation", strconv.Itoa(int(userId)))
	if conversation.HasResponse() {
		threeMinutes := 180
		a.Cache.Set(cacheKey, conversation, threeMinutes)
	} else {
		a.Cache.Delete(cacheKey)
	}

	return response
}
