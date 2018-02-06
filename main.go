package main

import (
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"os"

	"github.com/fmitra/dennis/postgres"
	"github.com/fmitra/dennis/expenses"
	"github.com/fmitra/dennis/wit"
	"github.com/fmitra/dennis/telegram"
	"github.com/fmitra/dennis/alphapoint"
)

type Config struct {
	Database struct {
		Host string `json:"host"`
		Port int32 `json:"port"`
		User string `json:"user"`
		Password string `json:"password"`
		Name string `json:"name"`
		SSLMode string `json:"ssl_mode"`
	} `json:"database"`
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

var webhookPath = fmt.Sprintf("/%s", telegram.Client.Token)
var config Config

func main() {
	// Load config
	loadConfig()

	// Set up DB
	setupDb()

	// Set up Wit.ai
	setupWitAi()

	// Set up Telegram
	setupTelegram()

	// Set up AlphaPoint
	setupAlphapoint()

	// Set up endpoints
	http.HandleFunc("/healthcheck", healthcheck)
	http.HandleFunc(webhookPath, webhook)

	// Run server
	log.Printf("main: starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadConfig() {
	file := "config.json"
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}

func setupAlphapoint() {
	alphaPointToken := config.AlphaPoint.Token
	alphapoint.Init(alphaPointToken)
}

func setupTelegram() {
	telegramToken := config.Telegram.Token
	telegram.Init(telegramToken)
}

func setupWitAi() {
	witToken := config.Wit.Token
	wit.Init(witToken)
}

func setupDb() {
	db := postgres.Config{
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Name,
		config.Database.Password,
		config.Database.SSLMode,
	}

	db.Open()
	postgres.Db.AutoMigrate(&expenses.Expense{})
}

func healthcheck(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok"))
}

func webhook(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	log.Printf("main: incoming message - %s", body)
	converse(body)
	w.Write([]byte("received"))
}
