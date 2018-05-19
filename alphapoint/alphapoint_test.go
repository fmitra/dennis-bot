package alphapoint

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeTestServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))
}

func TestAlphapoint(t *testing.T) {
	t.Run("Returns client with default config", func(t *testing.T) {
		token := "alphapointToken"
		alphapoint := Client(token)

		assert.Equal(t, "alphapointToken", alphapoint.Token)
		assert.Equal(t, BaseUrl, alphapoint.BaseUrl)
	})

	t.Run("Converts currency", func(t *testing.T) {
		response := `{
			"Realtime Currency Exchange Rate": {
				"5. Exchange Rate": ".7"
			}
		}`

		server := makeTestServer(response)
		defer server.Close()

		alphapoint := client{
			Token:   "alphapointToken",
			BaseUrl: server.URL,
		}

		fromISO := "USD"
		toISO := "SGD"
		forConversion := 20.00

		convertedAmount := alphapoint.Convert(fromISO, toISO, forConversion)
		assert.Equal(t, 14.0, convertedAmount)
	})
}
