package conversation

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/wit"
	mocks "github.com/fmitra/dennis-bot/test"
)

func TestGenericResponse(t *testing.T) {
	t.Run("Should return a list of possible responses", func(t *testing.T) {
		genericResponse := &GenericResponse{}
		assert.Equal(t, 1, len(genericResponse.GetResponses()))
	})

	t.Run("Should return a generic response", func(t *testing.T) {
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

		var incMessage telegram.IncomingMessage
		message := mocks.GetMockMessage("")
		json.Unmarshal(message, &incMessage)

		genericResponse := &GenericResponse{
			Context{
				Step:        0,
				WitResponse: witResponse,
				IncMessage:  incMessage,
			},
			&Actions{},
		}

		response, step := genericResponse.Respond()
		assert.Equal(t, -1, step)
		assert.Equal(t, BotResponse("This is a default message"), response)
	})
}
