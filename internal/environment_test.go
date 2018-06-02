package internal

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis/config"
	"github.com/fmitra/dennis/mocks"
	"github.com/fmitra/dennis/telegram"
	"github.com/fmitra/dennis/wit"
)

func TestEnvironment(t *testing.T) {
	t.Run("Should respond to healthcheck", func(t *testing.T) {
		configFile := "../config/config.json"
		config := config.LoadConfig(configFile)
		env := LoadEnv(config)

		req, err := http.NewRequest("GET", "/healthcheck", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(env.HealthCheck())

		handler.ServeHTTP(rr, req)
		assert.Equal(t, "ok", rr.Body.String())
	})

	t.Run("Should respond to webhook", func(t *testing.T) {
		configFile := "../config/config.json"
		config := config.LoadConfig(configFile)
		message := mocks.GetMockMessage()

		witResponse := `{
			"entities": {
				"amount": [],
				"datetime": [],
				"description": []
			}
		}`
		telegramServer := mocks.MakeTestServer("")
		witServer := mocks.MakeTestServer(witResponse)

		telegram.BaseUrl = fmt.Sprintf("%s/", telegramServer.URL)
		wit.BaseUrl = witServer.URL

		env := LoadEnv(config)

		req, err := http.NewRequest("POST", "/webook", bytes.NewBuffer(message))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(env.Webhook())

		handler.ServeHTTP(rr, req)
		assert.Equal(t, "received", rr.Body.String())
	})

	t.Run("Should load environment from config", func(t *testing.T) {
		configFile := "../config/config.json"
		config := config.LoadConfig(configFile)
		env := LoadEnv(config)

		assert.NotNil(t, env.db)
		assert.NotNil(t, env.cache)
		assert.NotNil(t, env.config)
		assert.NotNil(t, env.telegram)
	})
}
