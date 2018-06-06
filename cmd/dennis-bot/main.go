package main

import (
	"github.com/fmitra/dennis-bot/config"
	"github.com/fmitra/dennis-bot/internal"
)

func main() {
	// Set up the environment
	configFile := "config/config.json"
	env := internal.LoadEnv(config.LoadConfig(configFile))

	env.Start()
}
