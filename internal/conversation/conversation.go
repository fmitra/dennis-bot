package conversation

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/fmitra/dennis-bot/pkg/sessions"
	t "github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/wit"
)

const (
	ONBOARD_USER_INTENT      = "onboard_user_intent"
	TRACK_EXPENSE_INTENT     = "track_expense_intent"
	GET_EXPENSE_TOTAL_INTENT = "get_expense_total_intent"
)

// Intents describe what a user wants to accomplish from a conversation.
// They are responsible for the business logic behind any action we
// take in response to a user's incoming message and return a response
// for that action.
type Intent interface {
	// Determine's a user's intent behind an incoming message
	Respond() (BotResponse, int)

	// Returns a list of functions to deliver responses to a user.
	// Functions are listed in order of delivery from top to bottom.
	GetResponses() []func() (BotResponse, error)
}

// Context provides necessary info and methods for an Intent to process a response
type Context struct {
	Step        int
	WitResponse wit.WitResponse
	IncMessage  t.IncomingMessage
}

// Represents a user's place in a conversation, for example a conversation
// may have just recently been initialized or may be ongoing.
type Conversation struct {
	Intent string
	UserId int
	Step   int
}

// Helper method to ensure all Intents embeding context has access to the same
// response logic. Intent response methods should return a BotResponse and an error.
// Nil errors and empty responses are typically performed in validation steps (ex.
// ask a question and check for an answer). Valid responses increment a step,
// allowing a user to receive the next response in defined in the `GetResponses` array
func (cx *Context) Process(responses []func() (BotResponse, error)) (BotResponse, int) {
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

	return response, cx.Step
}

// Ends a conversation
func (cx *Context) EndConversation() {
	cx.Step = -1
}

func (c *Conversation) HasResponse() bool {
	return c.Step > -1
}

// Creates a new intent with additional context fields in order to
// determine a response
func (c *Conversation) GetIntent(w wit.WitResponse, inc t.IncomingMessage, a *Actions) Intent {
	context := Context{
		WitResponse: w,
		IncMessage:  inc,
		Step:        c.Step,
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

	intent := c.GetIntent(w, inc, a)
	response, nextStep := intent.Respond()
	c.Step = nextStep

	return response
}

func InferIntent(w wit.WitResponse) string {
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

func NewConversation(userId int, w wit.WitResponse) Conversation {
	intent := InferIntent(w)
	conversation := Conversation{
		Intent: intent,
		UserId: userId,
	}

	return conversation
}

func GetConversation(userId int, cache sessions.Session) (Conversation, error) {
	var conversation Conversation
	cacheKey := fmt.Sprintf("%s_conversation", strconv.Itoa(userId))
	cache.Get(cacheKey, &conversation)
	if conversation.UserId == 0 {
		return conversation, errors.New("No conversation found")
	}
	return conversation, nil
}

func GetResponse(w wit.WitResponse, inc t.IncomingMessage, a *Actions) BotResponse {
	userId := inc.GetUser().Id

	conversation, err := GetConversation(userId, a.Cache)
	if err != nil {
		conversation = NewConversation(userId, w)
	}

	response := conversation.Respond(w, inc, a)

	if conversation.HasResponse() {
		cacheKey := fmt.Sprintf("%s_conversation", strconv.Itoa(userId))
		a.Cache.Set(cacheKey, conversation)
	}

	return response
}
