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

func (t Telegram) sendUrl() (url string) {
	return fmt.Sprintf("%s%s/sendMessage", baseUrl, t.Token)
}

// Sends a message to Telegram. Sending a message
// requires the ID of a chat log (received from an IncomingMessage)
// and the text we are sending to the user.
func (t Telegram) send(chatId int, message string) {
	url := t.sendUrl()
	outMessage := OutgoingMessage{chatId, message}
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
	log.Printf("telegram: outgoing message - %s - %s", payload, body)
}
