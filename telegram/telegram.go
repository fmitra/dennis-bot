package telegram

import (
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"bytes"
)

var baseUrl = "https://api.telegram.org/bot"
var Client Telegram

type Telegram struct {
	Token string
}

// Set up client to run with Telegram token
func Init(token string) {
	Client = Telegram{token}
}

// Sends a message to Telegram. Sending a message
// requires the ID of a chat log (received from an IncomingMessage)
// and the text we are sending to the user.
func (t Telegram) Send(chatId int, message string) {
	url := fmt.Sprintf("%s%s/sendMessage", baseUrl, t.Token)
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
