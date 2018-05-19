package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("It returns AppConfig from JSON", func(t *testing.T) {
		file := "config.json"
		appConfig := LoadConfig(file)
		assert.IsType(t, AppConfig{}, appConfig)
	})
}
