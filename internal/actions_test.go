package internal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis-bot/config"
	"github.com/fmitra/dennis-bot/pkg/alphapoint"
	"github.com/fmitra/dennis-bot/pkg/expenses"
	"github.com/fmitra/dennis-bot/pkg/wit"
	mocks "github.com/fmitra/dennis-bot/test"
)

func TestActions(t *testing.T) {
	t.Run("Creates a new expense", func(t *testing.T) {
		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		env.cache.Delete("SGD_USD")
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
		alphapointResponse := `{
			"Realtime Currency Exchange Rate": {
				"5. Exchange Rate": ".7"
			}
		}`
		alphapointServer := mocks.MakeTestServer(alphapointResponse)
		defer alphapointServer.Close()

		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)

		alphapoint.BaseUrl = alphapointServer.URL

		a := &Actions{
			env:         env,
			witResponse: witResponse,
			userId:      mocks.TestUserId,
		}

		isCreated := a.createNewExpense()
		assert.True(t, isCreated)
	})

	t.Run("Creates a new expense from cache", func(t *testing.T) {
		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		env.cache.Delete("SGD_USD")
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
		alphapointResponse := `{
			"Realtime Currency Exchange Rate": {
				"5. Exchange Rate": ".7"
			}
		}`
		alphapointServer := mocks.MakeTestServer(alphapointResponse)

		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)

		alphapoint.BaseUrl = alphapointServer.URL

		a := &Actions{
			env:         env,
			witResponse: witResponse,
			userId:      mocks.TestUserId,
		}

		// Initial call without cache
		a.createNewExpense()

		// Second call should not hit server
		alphapointServer.Close()
		isCreated := a.createNewExpense()
		assert.True(t, isCreated)
	})

	t.Run("Should handle tracking intent from Wit.ai", func(t *testing.T) {
		var witResponse wit.WitResponse
		witResponseRaw := []byte(`{
			"entities": {
				"amount": [
					{ "value": "20 USD", "confidence": 100.00 }
				],
				"datetime": [
					{ "value": "", "confidence": 100.00 }
				],
				"description": [
					{ "value": "Food", "confidence": 100.00 }
				]
			}
		}`)
		json.Unmarshal(witResponseRaw, &witResponse)

		alphapointResponse := `{
			"Realtime Currency Exchange Rate": {
				"5. Exchange Rate": ".7"
			}
		}`
		alphapointServer := mocks.MakeTestServer(alphapointResponse)
		defer alphapointServer.Close()

		MessageMap = mocks.MessageMapMock
		alphapoint.BaseUrl = alphapointServer.URL
		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		a := &Actions{
			env:         env,
			userId:      mocks.TestUserId,
			witResponse: witResponse,
		}

		response := a.ProcessIntent()
		assert.Equal(t, BotResponse("Roger that!"), response)
	})

	t.Run("Should handle period query intent from Wit.ai", func(t *testing.T) {
		var witResponse wit.WitResponse
		witResponseRaw := []byte(`{
			"entities": {
				"amount": [],
				"datetime": [],
				"description": [],
				"total_spent": [
					{ "value": "month", "confidence": 100.00 }
				]
			}
		}`)
		json.Unmarshal(witResponseRaw, &witResponse)

		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		a := &Actions{
			env:         env,
			witResponse: witResponse,
			userId:      mocks.TestUserId,
		}
		a.env.db.Where("user_id = ?", mocks.TestUserId).
			Unscoped().
			Delete(expenses.Expense{})

		response := a.ProcessIntent()
		assert.Equal(t, BotResponse("You spent 0.00"), response)
	})

	t.Run("Should return a default message", func(t *testing.T) {
		var witResponse wit.WitResponse
		witResponseRaw := []byte(`{
			"entities": {
				"amount": [],
				"datetime": [],
				"description": [],
				"total_spent": []
			}
		}`)
		json.Unmarshal(witResponseRaw, &witResponse)
		MessageMap = mocks.MessageMapMock

		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		a := &Actions{
			env:         env,
			witResponse: witResponse,
			userId:      mocks.TestUserId,
		}

		response := a.ProcessIntent()
		assert.Equal(t, BotResponse("This is a default message"), response)
	})

	t.Run("Should return a error message", func(t *testing.T) {
		var witResponse wit.WitResponse
		witResponseRaw := []byte(`{
			"entities": {
				"amount": [],
				"datetime": [],
				"description": [],
				"total_spent": [
					{ "value": "foo", "confidence": 100.00 }
				]
			}
		}`)
		json.Unmarshal(witResponseRaw, &witResponse)
		MessageMap = mocks.MessageMapMock

		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		a := &Actions{
			env:         env,
			witResponse: witResponse,
			userId:      mocks.TestUserId,
		}

		response := a.ProcessIntent()
		assert.Equal(t, BotResponse("Whoops!"), response)
	})

	t.Run("It should return a historical total by period", func(t *testing.T) {
		rawWitResponse := []byte(`{
			"entities": {
				"amount": [],
				"datetime": [],
				"description": [],
				"total_spent": [
					{ "value": "month", "confidence": 100.00 }
				]
			}
		}`)
		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)
		MessageMap = mocks.MessageMapMock

		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		a := &Actions{
			env:         env,
			witResponse: witResponse,
			userId:      mocks.TestUserId,
		}
		a.env.db.Where("user_id = ?", mocks.TestUserId).
			Unscoped().
			Delete(expenses.Expense{})

		response := a.GetExpenseTotal()
		assert.Equal(t, BotResponse("You spent 0.00"), response)
	})

	t.Run("It should return error message for invalid period", func(t *testing.T) {
		rawWitResponse := []byte(`{
			"entities": {
				"amount": [],
				"datetime": [],
				"description": [],
				"total_spent": [
					{ "value": "foo", "confidence": 100.00 }
				]
			}
		}`)
		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)
		MessageMap = mocks.MessageMapMock

		configFile := "../config/config.json"
		env := LoadEnv(config.LoadConfig(configFile))
		a := &Actions{
			env:         env,
			witResponse: witResponse,
			userId:      mocks.TestUserId,
		}
		a.env.db.Where("user_id = ?", mocks.TestUserId).
			Unscoped().
			Delete(expenses.Expense{})

		response := a.GetExpenseTotal()
		assert.Equal(t, BotResponse("Whoops!"), response)
	})
}
