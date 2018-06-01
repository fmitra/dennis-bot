package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var BaseUrl = "https://api.telegram.org/bot"

type Telegram interface {
	SetWebhook() int
	Send(chatId int, message string) int
}

type Client struct {
	Token   string
	Domain  string
	BaseUrl string
}

// Convenience function to return an API client with default
// base URL
func NewClient(token string, domain string) *Client {
	return &Client{
		Token:   token,
		Domain:  domain,
		BaseUrl: BaseUrl,
	}
}

// Sets the bot domain webhook with Telegram
func (c Client) SetWebhook() int {
	webhook := fmt.Sprintf("%s/%s", c.Domain, c.Token)
	url := fmt.Sprintf("%s%s/setWebhook?url=%s", c.BaseUrl, c.Token, webhook)
	resp, httpErr := http.Get(url)
	defer resp.Body.Close()

	if httpErr != nil {
		panic(httpErr)
	}

	return resp.StatusCode
}

// Sends a message to Telegram. Sending a message
// requires the ID of a chat log (received from an IncomingMessage)
// and the text we are sending to the user.
func (c Client) Send(chatId int, message string) int {
	url := fmt.Sprintf("%s%s/sendMessage", c.BaseUrl, c.Token)
	contentType := "application/json"
	outMessage := OutgoingMessage{chatId, message}
	payload, _ := json.Marshal(outMessage)
	resp, err := http.Post(url, contentType, bytes.NewReader(payload))
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("telegram: outgoing message - %s - %s", payload, body)

	return resp.StatusCode
}
