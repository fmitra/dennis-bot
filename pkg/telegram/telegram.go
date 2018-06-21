// Package telegram implements a wrapper for the Telegram Bot API.
// Telegram is the chat platform the bot service runs on.
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

// BaseURL for Telegram
var BaseURL = "https://api.telegram.org/bot"

// Telegram is an interface to provide utility methods to interact with
// the Telegram API.
type Telegram interface {
	SetWebhook() int
	Send(chatID int, message string) int
	SendAction(chatID int, action string) int
}

// Client is a consumer of the Telegram API.
type Client struct {
	Token   string
	Domain  string
	BaseURL string
}

// NewClient returns a Client with default BaseUrl to interact with the Telegram API.
func NewClient(token string, domain string) *Client {
	return &Client{
		Token:   token,
		Domain:  domain,
		BaseURL: BaseURL,
	}
}

// SetWebhook update's Telegram with the location of the bot webhook.
// Failure at this step causes a panic as the bot cannot run if it Telegram
// is not set up with the webhook. Return's an HTTP status code.
func (c *Client) SetWebhook() int {
	webhook := fmt.Sprintf("%s/%s", c.Domain, c.Token)
	url := fmt.Sprintf("%s%s/setWebhook?url=%s", c.BaseURL, c.Token, webhook)
	resp, err := http.Get(url)
	if err != nil {
		log.Panicf("telegram: unable to set webhook - %s", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

// Send sends a text message to a User. Returns an HTTP status code.
func (c *Client) Send(chatID int, message string) int {
	url := fmt.Sprintf("%s%s/sendMessage", c.BaseURL, c.Token)
	contentType := "application/json"
	outMessage := OutgoingMessage{chatID, message}
	payload, err := json.Marshal(outMessage)
	if err != nil {
		log.Printf("telegram: cannot send invalid message format")
		errorCode := 400
		return errorCode
	}

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

	err = retry.Retry(
		request,
		strategy.Limit(10),
		strategy.Backoff(backoff.Exponential(time.Second, 2)),
	)
	if err != nil {
		log.Printf("telegram: failed to send message - %s", err)
	}

	respBody.Close()
	return statusCode
}

// SendAction sends a ChatAction to a User. This is used to alert
// the user that we have received their message and will respond soon.
// The most commong usage is to send a typing indicator. Returns an
// HTTP status code.
func (c *Client) SendAction(chatID int, action string) int {
	url := fmt.Sprintf("%s%s/sendChatAction", c.BaseURL, c.Token)
	contentType := "application/json"
	chatAction := ChatAction{chatID, action}
	payload, err := json.Marshal(chatAction)
	errorCode := 400
	if err != nil {
		log.Printf("telegram: cannot send invalid message format - %s", err)
		return errorCode
	}

	var statusCode int
	resp, err := http.Post(url, contentType, bytes.NewReader(payload))
	if err != nil {
		log.Printf("telegram: failed to send action - %s", err)
		return errorCode
	}

	statusCode = resp.StatusCode
	resp.Body.Close()
	return statusCode
}
