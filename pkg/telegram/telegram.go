package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
)

var BaseUrl = "https://api.telegram.org/bot"

type Telegram interface {
	SetWebhook() int
	Send(chatId int, message string) int
	SendAction(chatId int, action string) int
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
func (c *Client) SetWebhook() int {
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
func (c *Client) Send(chatId int, message string) int {
	url := fmt.Sprintf("%s%s/sendMessage", c.BaseUrl, c.Token)
	contentType := "application/json"
	outMessage := OutgoingMessage{chatId, message}
	payload, _ := json.Marshal(outMessage)

	var respBody io.ReadCloser
	var statusCode int
	request := func(attempt uint) error {
		resp, err := http.Post(url, contentType, bytes.NewReader(payload))
		respBody = resp.Body
		statusCode = resp.StatusCode

		if err != nil {
			return err
		}

		return nil
	}

	err := retry.Retry(
		request,
		strategy.Limit(10),
		strategy.Backoff(backoff.Exponential(time.Second, 2)),
	)
	if err != nil {
		log.Printf("telegram: failed to send message")
	}

	respBody.Close()
	return statusCode
}

// Sends a processing notice to the user. Used to alert the user
// that the bot has received their message and will respond soon.
// Most common usage is a typing indicator.
func (c *Client) SendAction(chatId int, action string) int {
	url := fmt.Sprintf("%s%s/sendChatAction", c.BaseUrl, c.Token)
	contentType := "application/json"
	chatAction := ChatAction{chatId, action}
	payload, _ := json.Marshal(chatAction)

	var statusCode int
	resp, err := http.Post(url, contentType, bytes.NewReader(payload))
	if err != nil {
		log.Printf("telegram: failed to send action")
		statusCode = 400
	}

	statusCode = resp.StatusCode
	resp.Body.Close()
	return statusCode
}
