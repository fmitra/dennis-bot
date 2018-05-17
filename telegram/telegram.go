package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const BaseUrl = "https://api.telegram.org/bot"

type HttpLib interface {
	Get(url string) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type client struct {
	Token string
	Domain string
	BaseUrl string
	Http    HttpLib
}

func Client(token string, domain string, httpLib HttpLib) *client {
	return &client{
		Token: token,
		Domain: domain,
		BaseUrl: BaseUrl,
		Http: httpLib,
	}
}

// Sets the bot domain webhook with Telegram
func (c client) SetWebhook() {
	webhook := fmt.Sprintf("%s/%s", c.Domain, c.Token)
	url := fmt.Sprintf("%s%s/setWebhook?url=%s", c.BaseUrl, c.Token, webhook)
	resp, err := c.Http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}
}

// Sends a message to Telegram. Sending a message
// requires the ID of a chat log (received from an IncomingMessage)
// and the text we are sending to the user.
func (c client) Send(chatId int, message string) {
	url := fmt.Sprintf("%s%s/sendMessage", c.BaseUrl, c.Token)
	contentType := "application/json"
	outMessage := OutgoingMessage{chatId, message}
	payload, _ := json.Marshal(outMessage)
	resp, err := c.Http.Post(url, contentType, bytes.NewReader(payload))
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("telegram: outgoing message - %s - %s", payload, body)
}
