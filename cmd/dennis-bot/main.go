package main

import (
	"os"

	"github.com/fmitra/dennis-bot/config"
	"github.com/fmitra/dennis-bot/internal"
)

func main() {
	configFile := os.Getenv("DENNIS_BOT_CONFIG")
	if configFile == "" {
		configFile = "config/config.json"
	}

	env := internal.LoadEnv(config.LoadConfig(configFile))

	// Set up the dependencies and start the HTTP handlers
	env.Start()
}
