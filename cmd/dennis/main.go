package main

import (
	"github.com/fmitra/dennis/config"
	"github.com/fmitra/dennis/internal"
)

func main() {
	// Set up the environment
	configFile := "config/config.json"
	env := internal.LoadEnv(config.LoadConfig(configFile))

	env.Start()
}
