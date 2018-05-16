package config

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	Database struct {
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

func LoadConfig() AppConfig {
	var config AppConfig

	file := "config/config.json"
	configFile, err := os.Open(file)
	defer configFile.Close()

	if err != nil {
		panic(err)
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}
