package main

import (
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"bytes"
	"os"
)

var baseUrl = "https://api.telegram.org/bot"
var telegram = Telegram{os.Getenv("TELEGRAM_AUTH_TOKEN")}

type Telegram struct {
	Token string
}

// ref: https://core.telegram.org/bots/api#update
// Represents a Telegram Update object. This payload is sent
// to the webhook whenever a user messages us. Message field
// is optional but for our use case this all we care about at
// the moment
type IncomingMessage struct {
	UpdateId int `json:"update_id"`
	Message struct {
		MessageId int `json:"message_id"`
		Date int `json:"date"`
		Text string `json:"text"`
		From struct {
			Id int `json:"id"`
			FirstName string `json:"first_name"`
			LastName string `json:"last_name"`
			UserName string `json:"username"`
		} `json:"from"`
		Chat struct {
			Id int `json:"id"`
			FirstName string `json:"first_name"`
			LastName string `json:"last_name"`
			UserName string `json:"username"`
		} `json:"chat"`
	} `json:"message"`
}

// Bot response to an `IncomingMessage`
type OutgoingMessage struct {
	ChatId int `json:"chat_id"`
	Text string `json:"text"`
}

// Receives incoming messages from Telegram and processes a response
func (t Telegram) Webhook(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	log.Printf("bot: incoming message - %s", body)
	parseMessage(body)
	defer req.Body.Close()
	w.Write([]byte("received"))
}

func (t Telegram) sendUrl() (string) {
	return fmt.Sprintf("%s%s/sendMessage", baseUrl, telegram.Token)
}

// Parses an incoming message from Telegram. If the message is valid,
// we send a response in a new thread
func parseMessage(payload []byte) {
	var incMessage IncomingMessage
	err := json.Unmarshal(payload, &incMessage)
	if err != nil {
		log.Printf("bot: invalid payload - %s", err)
		return
	}

	go sendMessage(incMessage)
}

// Sends a message to Telegram
func sendMessage(incMessage IncomingMessage) {
	url := telegram.sendUrl()
	outMessage := OutgoingMessage{incMessage.Message.Chat.Id, "Hi I'm Dennis"}
	payload, _ := json.Marshal(outMessage)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("bot: outgoing message - %s - %s", payload, body)
}
