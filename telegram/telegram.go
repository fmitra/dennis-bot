package telegram

import (
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"io"
	"fmt"
	"bytes"
)

var Client Telegram
const baseUrl = "https://api.telegram.org/bot"

type HttpLib interface {
	Get(url string) (resp http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp http.Response, err error)
}

type Telegram struct {
	Token string
	Domain string
	BaseUrl string
	Http HttpLib
}

// Set up client to run with Telegram token
func Init(token string, domain string, httpLib HttpLib) chan bool {
	Client = Telegram{
		token,
		domain,
		baseUrl,
		httpLib,
	}
	channel := make(chan bool)

	go func() {
		channel <- true
		Client.SetWebhook()
	}()

	return channel
}

// Sets the bot domain webhook with Telegram
func (t Telegram) SetWebhook() {
	webhook := fmt.Sprintf("%s/%s", t.Domain, t.Token)
	url := fmt.Sprintf("%s%s/setWebhook?url=%s", t.BaseUrl, t.Token, webhook)
	resp, err := t.Http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}
}

// Sends a message to Telegram. Sending a message
// requires the ID of a chat log (received from an IncomingMessage)
// and the text we are sending to the user.
func (t Telegram) Send(chatId int, message string) {
	url := fmt.Sprintf("%s%s/sendMessage", t.BaseUrl, t.Token)
	contentType := "application/json"
	outMessage := OutgoingMessage{chatId, message}
	payload, _ := json.Marshal(outMessage)
	resp, err := t.Http.Post(url, contentType, bytes.NewReader(payload))
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("telegram: outgoing message - %s - %s", payload, body)
}
