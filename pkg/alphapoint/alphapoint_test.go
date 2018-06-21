package alphapoint

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis-bot/test"
)

func TestAlphapoint(t *testing.T) {
	t.Run("Returns client with default config", func(t *testing.T) {
		token := "alphapointToken"
		alphapoint := NewClient(token)

		assert.Equal(t, "alphapointToken", alphapoint.Token)
		assert.Equal(t, BaseURL, alphapoint.BaseURL)
	})

	t.Run("Converts currency", func(t *testing.T) {
		response := `{
			"Realtime Currency Exchange Rate": {
				"5. Exchange Rate": ".7"
			}
		}`

		server := mocks.MakeTestServer(response)
		defer server.Close()

		alphapoint := Client{
			Token:   "alphapointToken",
			BaseURL: server.URL,
		}

		fromISO := "USD"
		toISO := "SGD"
		forConversion := 20.00

		convertedAmount, conversion := alphapoint.Convert(
			fromISO,
			toISO,
			forConversion,
		)

		expectedConversion := &Conversion{
			fromISO,
			toISO,
			0.7,
		}

		assert.Equal(t, expectedConversion, conversion)
		assert.Equal(t, 14.0, convertedAmount)
	})
}
