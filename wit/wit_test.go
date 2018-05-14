package wit

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

func TestWit(t *testing.T) {
	t.Run("Returns WitResponse", func(t *testing.T) {
		response := `{
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
		server := makeTestServer(response)
		defer server.Close()

		witAi := WitAi{
			Token: "witAiToken",
			BaseUrl: server.URL,
			ApiVersion: "20180128",
		}

		witResponse := witAi.ParseMessage("Hello world")
		assert.IsType(t, WitResponse{}, witResponse)
	})

	t.Run("Fails if JSON is not returned", func(t *testing.T) {
		server := makeTestServer(`{not valid json}`)
		defer server.Close()

		witAi := WitAi{
			Token: "witAiToken",
			BaseUrl: server.URL,
			ApiVersion: "20180128",
		}

		assert.Panics(t, func() { witAi.ParseMessage("Hello world") })
	})
}
