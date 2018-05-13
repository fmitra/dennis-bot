package alphapoint

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func makeTestServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))
}

func TestAlphapoint(t *testing.T) {
	t.Run("Sets Client on init", func(t *testing.T) {
		token := "alphapointToken"
		httpLib := &http.Client{}
		Init(token, httpLib)

		assert.Equal(t, "alphapointToken", Client.Token)
		assert.Equal(t, "https://www.alphavantage.co/query", Client.BaseUrl)
		assert.Equal(t, httpLib, Client.Http)
	})

	t.Run("Converts currency", func(t *testing.T) {
		response := `{
			"Realtime Currency Exchange Rate": {
				"5. Exchange Rate": ".7"
			}
		}`

		server := makeTestServer(response)
		defer server.Close()

		alphapoint := AlphaPoint{
			Token: "alphapointToken",
			BaseUrl: server.URL,
			Http: &http.Client{},
		}

		fromISO := "USD"
		toISO := "SGD"
		forConversion := 20.00

		convertedAmount := alphapoint.Convert(fromISO, toISO, forConversion)
		assert.Equal(t, 14.0, convertedAmount)
	})
}
