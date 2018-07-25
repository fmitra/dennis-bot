// Package telegram implements a wrapper for the Telegram Bot API.
// Telegram is the chat platform the bot service runs on.
package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
)

// BaseURL for Telegram
const BaseURL = "https://api.telegram.org/bot"

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

// SendDocument sends a file to a User. After a file is delivered,
// it is deleted from the bot.
func (c *Client) SendDocument(chatID int, fileName string) int {
	defer os.Remove(fileName)
	url := fmt.Sprintf("%s%s/sendDocument", c.BaseURL, c.Token)
	errorCode := 400

	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("telegram: cannot open file - %s", err)
		return errorCode
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("document", fileName)
	if err != nil {
		log.Printf("telegram: cannot create form file - %s", err)
		return errorCode
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Printf("telegram: cannot copy file - %s", err)
		return errorCode
	}

	writer.WriteField("chat_id", strconv.Itoa(int(chatID)))
	err = writer.Close()
	if err != nil {
		log.Printf("telegram: writer close error - %s", err)
		return errorCode
	}

	// TODO Exponenial retry?
	var statusCode int
	resp, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		log.Printf("telegram: failed to send document - %s", err)
		return errorCode
	}

	statusCode = resp.StatusCode
	resp.Body.Close()
	return statusCode
}
