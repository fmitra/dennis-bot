package main

import (
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fmitra/dennis/postgres"
	"github.com/fmitra/dennis/expenses"
	"github.com/fmitra/dennis/wit"
	"github.com/fmitra/dennis/telegram"
	"github.com/fmitra/dennis/alphapoint"
)

var webhookPath = fmt.Sprintf("/%s", telegram.Client.Token)

func main() {
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

func setupAlphapoint() {
	alphaPointToken := os.Getenv("ALPHAPOINT_AUTH_TOKEN")
	alphapoint.Init(alphaPointToken)
}

func setupTelegram() {
	telegramToken := os.Getenv("TELEGRAM_AUTH_TOKEN")
	telegram.Init(telegramToken)
}

func setupWitAi() {
	witToken := os.Getenv("WITAI_AUTH_TOKEN")
	wit.Init(witToken)
}

func setupDb() {
	config := postgres.Config{
		"0.0.0.0",
		5432,
		"dennis",
		"dennis_test",
		"dennis",
		"disable",
	}

	config.Open()
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
