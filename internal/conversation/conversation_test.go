package conversation

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis-bot/config"
	"github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/wit"
	mocks "github.com/fmitra/dennis-bot/test"
)

func TestConversation(t *testing.T) {
	t.Run("Should return boolean for response check", func(t *testing.T) {
		conversation := &Conversation{}
		hasResponse := conversation.HasResponse()
		assert.True(t, hasResponse)

		conversation = &Conversation{
			Step: -1,
		}
		hasResponse = conversation.HasResponse()
		assert.False(t, hasResponse)
	})

	t.Run("Should create a new conversation", func(t *testing.T) {
		rawWitResponse := []byte(`{
			"entities": {
				"amount": [
					{ "value": "20 SGD", "confidence": 100.00 }
				],
				"datetime": [
					{ "value": "", "confidence": 100.00 }
				],
				"description": [
					{ "value": "Food", "confidence": 100.00 }
				]
			}
		}`)
		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)

		conversation := NewConversation(mocks.TestUserId, witResponse)
		assert.Equal(t, mocks.TestUserId, conversation.UserId)
		assert.Equal(t, TRACK_EXPENSE_INTENT, conversation.Intent)
	})

	t.Run("Should infer user intent from WitResponse", func(t *testing.T) {
		var rawWitResponse []byte
		var witResponse wit.WitResponse

		rawWitResponse = []byte(`{
			"entities": {
				"amount": [
					{ "value": "20 SGD", "confidence": 100.00 }
				],
				"datetime": [
					{ "value": "", "confidence": 100.00 }
				],
				"description": [
					{ "value": "Food", "confidence": 100.00 }
				]
			}
		}`)
		json.Unmarshal(rawWitResponse, &witResponse)
		assert.Equal(t, TRACK_EXPENSE_INTENT, InferIntent(witResponse))

		rawWitResponse = []byte(`{
			"entities": {
				"amount": [
					{ "value": "20 SGD", "confidence": 100.00 }
				],
				"datetime": [
					{ "value": "", "confidence": 100.00 }
				],
				"description": []
			}
		}`)
		json.Unmarshal(rawWitResponse, &witResponse)
		assert.Equal(t, TRACK_EXPENSE_INTENT, InferIntent(witResponse))

		rawWitResponse = []byte(`{
			"entities": {
				"amount": [],
				"datetime": [],
				"description": [],
				"total_spent": [
					{ "value": "month", "confidence": 100.00 }
				]
			}
		}`)
		json.Unmarshal(rawWitResponse, &witResponse)
		assert.Equal(t, GET_EXPENSE_TOTAL_INTENT, InferIntent(witResponse))

		rawWitResponse = []byte(`{
			"entities": {
				"amount": [],
				"datetime": [],
				"description": [],
				"total_spent": []
			}
		}`)
		json.Unmarshal(rawWitResponse, &witResponse)
		assert.Equal(t, "", InferIntent(witResponse))
	})

	t.Run("Should get conversation from cache", func(t *testing.T) {
		configFile := "../../config/config.json"
		appConfig := config.LoadConfig(configFile)
		cache := GetSession(appConfig)
		userId := 124
		conversation := Conversation{
			Intent: ONBOARD_USER_INTENT,
			UserId: userId,
		}
		cacheKey := "124_conversation"

		cache.Set(cacheKey, conversation)
		cachedConversation, err := GetConversation(userId, cache)

		assert.NoError(t, err)
		assert.Equal(t, cachedConversation, conversation)
	})

	t.Run("Should return error when fetching conversation from cache", func(t *testing.T) {
		configFile := "../../config/config.json"
		appConfig := config.LoadConfig(configFile)
		cache := GetSession(appConfig)
		userId := 125

		cachedConversation, err := GetConversation(userId, cache)
		assert.EqualError(t, err, "No conversation found")
		assert.Equal(t, cachedConversation, Conversation{})
	})

	t.Run("Should retrieve an Intent", func(t *testing.T) {
		conversation := &Conversation{
			Intent: ONBOARD_USER_INTENT,
		}
		witResponse := wit.WitResponse{}
		incMessage := telegram.IncomingMessage{}
		actions := &Actions{}
		intent := conversation.GetIntent(witResponse, incMessage, actions)
		assert.IsType(t, &OnboardUser{}, intent)
	})

	t.Run("Should get response from Intent in correct order", func(t *testing.T) {
		MessageMap = mocks.MessageMapMock
		conversation := &Conversation{
			Intent: ONBOARD_USER_INTENT,
		}
		actions := &Actions{}
		witResponse := wit.WitResponse{}
		incMessage := telegram.IncomingMessage{}
		message := mocks.GetMockMessage("")
		json.Unmarshal(message, &incMessage)

		// Starts at step 0
		assert.Equal(t, 0, conversation.Step)

		// First response requests password
		response := conversation.Respond(witResponse, incMessage, actions)
		assert.Equal(t, BotResponse("What's your password?"), response)
		assert.Equal(t, 1, conversation.Step)

		// Second response requests confirmation
		message = mocks.GetMockMessage("foo")
		json.Unmarshal(message, &incMessage)
		response = conversation.Respond(witResponse, incMessage, actions)
		assert.Equal(t, BotResponse("Your password is foo"), response)
		assert.Equal(t, 2, conversation.Step)

		// Invalid response prevents user from reaching step 3
		message = mocks.GetMockMessage("invalid answer")
		json.Unmarshal(message, &incMessage)
		response = conversation.Respond(witResponse, incMessage, actions)
		assert.Equal(t, BotResponse("I didn't understand that"), response)
		assert.Equal(t, 2, conversation.Step)

		// Answering no to password confirmation ends the conversation
		message = mocks.GetMockMessage("No")
		json.Unmarshal(message, &incMessage)
		response = conversation.Respond(witResponse, incMessage, actions)
		assert.Equal(t, BotResponse("Okay try again later"), response)
		assert.Equal(t, -1, conversation.Step)

		// After receiving a negative step, all future respones are empty
		message = mocks.GetMockMessage("Hello?")
		json.Unmarshal(message, &incMessage)
		response = conversation.Respond(witResponse, incMessage, actions)
		assert.Equal(t, BotResponse(""), response)
		assert.Equal(t, -1, conversation.Step)

		// Manually edit the step so we can continue the final tests
		conversation.Step = 2

		// When steps are iterated past the number of responses, we should
		// reset the step to -1 to end the conversation
		message = mocks.GetMockMessage("Yes")
		json.Unmarshal(message, &incMessage)
		response = conversation.Respond(witResponse, incMessage, actions)
		assert.Equal(t, BotResponse("Outro message"), response)
		assert.Equal(t, -1, conversation.Step)
	})
}
