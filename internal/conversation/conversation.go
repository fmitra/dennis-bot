// Package conversation manages the process of building a response for
// a user and taking action for any user request.
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
	// OnboardUserIntent is an intent designed to onboard a user into the bot
	// by creating a user account
	OnboardUserIntent = "onboard_user_intent"

	// TrackExpenseIntent is a intent designed to tracks a user's expense
	TrackExpenseIntent = "track_expense_intent"

	// GetExpenseTotalIntent is an intent to returns total expense
	// history for the user
	GetExpenseTotalIntent = "get_expense_total_intent"
)

// Intent describe the objective of a user. They are responsible for the
// business logic around any action we perform in response to a user's
// incoming message and contain all responses for that action.
type Intent interface {
	// Returns a list of functions to deliver responses to a user.
	// Functions are listed in order of delivery from top to bottom.
	GetResponses() []func() (BotResponse, error)
}

// Conversation represents a user's place in a conversation, for example a conversation
// may have just recently been initialized or may be ongoing. Conversations instantiate
// structs adhering to the Intent interface to build responses for the user.
type Conversation struct {
	IntentType  string            // String representing the the Intent of the user
	UserID      uint              // The user ID from the incoming chat service (ex. Telegram User ID)
	BotUserID   uint              // The ID of the user account (if it exists) in our system
	Step        int               // A user's place in a conversation
	AuxData     string            // Optional auxiliary info that we can set while processing a response
	WitResponse wit.Response      // Wit.ai API response to a user's message
	IncMessage  t.IncomingMessage // Raw user message as received by Telegram
}

// SetLastUserMessage sets the most recently received telegram message and Wit.ai
// parsing of the message to the Conversation struct.
func (c *Conversation) SetLastUserMessage(w wit.Response, inc t.IncomingMessage) {
	c.WitResponse = w
	c.IncMessage = inc
}

// ProcessResponses takes a list of methods that return a BotResponse and error
// to determine the next repsonse to deliver to a user. Nil errors and empty responses
// are typically performed in validation steps (ex. ask a question and check for an
// answer). Valid responses incremetn a step, allowing a user to receive the next
// response in the list.
func (c *Conversation) ProcessResponses(responses []func() (BotResponse, error)) BotResponse {
	response, err := responses[c.Step]()

	// We continue iterating through responses until we find the first
	// non-empty response.
	for response == BotResponse("") && err == nil {
		c.Step++
		response, err = responses[c.Step]()
	}

	conversationNotOver := err == nil && c.Step != -1
	if conversationNotOver {
		c.Step++
	}

	noResponsesAvailable := c.Step >= len(responses)
	if noResponsesAvailable {
		c.EndConversation()
	}

	return response
}

// SkipResponse returns an empty BotResponse and nil error. When processing a list
// of response functions, nil errors/empty responses are skipped.
func (c *Conversation) SkipResponse() (BotResponse, error) {
	return BotResponse(""), nil
}

// EndConversation sets the Covnerstation step to -1, indicating no further responses
// are available.
func (c *Conversation) EndConversation() {
	finalStep := -1
	c.Step = finalStep
}

// HasResponse checks if Conversation step > -1, indicating there are still responses
// to send to a User.
func (c *Conversation) HasResponse() bool {
	finalStep := -1
	return c.Step > finalStep
}

// GetIntent creates an a new Intent with the embedded Conversation to provide context.
// A Conversation create's an Intent in order to formulate a BotResponse.
func (c *Conversation) GetIntent(a *Actions) Intent {
	switch c.IntentType {
	case OnboardUserIntent:
		return &OnboardUser{c, a}
	case TrackExpenseIntent:
		return &TrackExpense{c, a}
	case GetExpenseTotalIntent:
		return &GetExpenseTotal{c, a}
	default:
		return &GenericResponse{c, a}
	}
}

// Respond returns a response to the user. Responses are controlled by Intent types
// which we create on demand using information from the surrounding Conversation.
func (c *Conversation) Respond(a *Actions) BotResponse {
	if !c.HasResponse() {
		return BotResponse("")
	}

	intent := c.GetIntent(a)
	responses := intent.GetResponses()
	return c.ProcessResponses(responses)
}

// InferIntent determines which Intent to instantiate for the Conversation.
func InferIntent(w wit.Response, botUserID uint) string {
	// We can only process a user's request if they have an account
	// in our system, otherwise we force them into an onboarding flow.
	noID := 0
	if int(botUserID) == noID {
		return OnboardUserIntent
	}

	overview := w.GetMessageOverview()
	switch overview {
	case wit.TrackingRequestedSuccess:
		return TrackExpenseIntent
	case wit.TrackingRequestedError:
		return TrackExpenseIntent
	case wit.ExpenseTotalRequestedSuccess:
		return GetExpenseTotalIntent
	default:
		return ""
	}
}

// NewConversation creates a new conversation between the bot and the user. If the user
// has an existing account, we associate their account ID with the user
// ID of their chat service.
func NewConversation(userID uint, w wit.Response, a *Actions) Conversation {
	manager := users.NewUserManager(a.Db)
	botUser := manager.GetByTelegramID(userID)

	intent := InferIntent(w, botUser.ID)
	conversation := Conversation{
		IntentType: intent,
		UserID:     userID,
		BotUserID:  botUser.ID,
	}

	return conversation
}

// GetConversation checks if an ongoing Conversation exists in the cache.
func GetConversation(userID uint, cache sessions.Session) (Conversation, error) {
	var conversation Conversation
	cacheKey := fmt.Sprintf("%s_conversation", strconv.Itoa(int(userID)))
	err := cache.Get(cacheKey, &conversation)
	if err != nil {
		return conversation, errors.New("no conversation found")
	}

	if !conversation.HasResponse() {
		return conversation, errors.New("no responses available")
	}
	return conversation, nil
}

// GetResponse creates or retrieves a Conversation in order to return the
// next available response.
func GetResponse(w wit.Response, inc t.IncomingMessage, a *Actions) BotResponse {
	userID := inc.GetUser().ID

	conversation, err := GetConversation(userID, a.Cache)
	if err != nil {
		conversation = NewConversation(userID, w, a)
	}

	conversation.SetLastUserMessage(w, inc)
	response := conversation.Respond(a)

	// Check if there are additional responses available. If responses are found,
	// we cache the conversation so the user can pick up where they left off
	cacheKey := fmt.Sprintf("%s_conversation", strconv.Itoa(int(userID)))
	if conversation.HasResponse() {
		threeMinutes := 180
		a.Cache.Set(cacheKey, conversation, threeMinutes)
	} else {
		a.Cache.Delete(cacheKey)
	}

	return response
}
