// Package config contains settings for the bot service and all related packages.
package config

import (
	"encoding/json"
	"os"
)

// AppConfig is a JSON config for the bot.
type AppConfig struct {
	SecretKey string `json:"secret_key"`
	Database  struct {
		Host     string `json:"host"`
		Port     int32  `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
		SSLMode  string `json:"ssl_mode"`
	} `json:"database"`
	Redis struct {
		Host     string `json:"host"`
		Port     int32  `json:"port"`
		Password string `json:"password"`
		Db       int    `json:"db"`
	} `json:"redis"`
	BotDomain  string `json:"bot_domain"`
	AlphaPoint struct {
		Token string `json:"token"`
	} `json:"alphapoint"`
	Telegram struct {
		Token string `json:"token"`
	} `json:"telegram"`
	Wit struct {
		Token string `json:"token"`
	} `json:"wit"`
}

// LoadConfig loads a JSON config from file (ex. config/config.json)
func LoadConfig(file string) AppConfig {
	var config AppConfig

	configFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}
