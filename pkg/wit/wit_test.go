package wit

import (
	"testing"

	mocks "github.com/fmitra/dennis-bot/test"
	"github.com/stretchr/testify/assert"
)

func TestWit(t *testing.T) {
	t.Run("Returns client with default config", func(t *testing.T) {
		witAi := NewClient("witAiToken")

		assert.Equal(t, BaseURL, witAi.BaseURL)
		assert.Equal(t, APIVersion, witAi.APIVersion)
	})

	t.Run("Returns Response", func(t *testing.T) {
		rawResponse := `{
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
		}`
		server := mocks.MakeTestServer(rawResponse)
		defer server.Close()

		witAi := Client{
			Token:      "witAiToken",
			BaseURL:    server.URL,
			APIVersion: "20180128",
		}

		response := witAi.ParseMessage("Hello world")
		assert.IsType(t, Response{}, response)
	})

	t.Run("Returns zero value Response on error", func(t *testing.T) {
		server := mocks.MakeTestServer(`{not valid json}`)
		defer server.Close()

		witAi := Client{
			Token:      "witAiToken",
			BaseURL:    server.URL,
			APIVersion: "20180128",
		}

		response := witAi.ParseMessage("Hello world")
		assert.Equal(t, Response{}, response)
	})
}
