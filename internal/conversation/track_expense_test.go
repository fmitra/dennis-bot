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

func TestTrackExpense(t *testing.T) {
	t.Run("Should return a list of possible responses", func(t *testing.T) {
		trackExpense := &TrackExpense{}
		assert.Equal(t, 1, len(trackExpense.GetResponses()))
	})

	t.Run("Should return success tracking message", func(t *testing.T) {
		MessageMap = mocks.MessageMapMock
		configFile := "../../config/config.json"
		appConfig := config.LoadConfig(configFile)
		db, _ := GetDb(appConfig)
		cache := GetSession(appConfig)
		action := &Actions{
			db,
			cache,
			appConfig,
		}

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

		trackExpense := &TrackExpense{
			Context{
				Step:        0,
				WitResponse: witResponse,
				IncMessage:  incMessage,
			},
			action,
		}
		response, step := trackExpense.Respond()
		assert.Equal(t, BotResponse("Roger that!"), response)
		assert.Equal(t, -1, step)
	})

	t.Run("Should return error tracking message", func(t *testing.T) {
		MessageMap = mocks.MessageMapMock
		configFile := "../../config/config.json"
		appConfig := config.LoadConfig(configFile)
		db, _ := GetDb(appConfig)
		cache := GetSession(appConfig)
		action := &Actions{
			db,
			cache,
			appConfig,
		}

		rawWitResponse := []byte(`{
			"entities": {
				"amount": [],
				"datetime": [
					{ "value": "", "confidence": 100.00 }
				],
				"description": []
			}
		}`)
		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)

		var incMessage telegram.IncomingMessage
		message := mocks.GetMockMessage("")
		json.Unmarshal(message, &incMessage)

		trackExpense := &TrackExpense{
			Context{
				Step:        0,
				WitResponse: witResponse,
				IncMessage:  incMessage,
			},
			action,
		}
		response, step := trackExpense.Respond()
		assert.Equal(t, BotResponse("Whoops!"), response)
		assert.Equal(t, -1, step)
	})
}
