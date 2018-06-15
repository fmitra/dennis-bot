package wit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
)

var (
	// Base URL for Wit.ai
	BaseUrl = "https://api.wit.ai"

	// Wit.ai API Version
	ApiVersion = "20180128"
)

type Wit interface {
	ParseMessage(message string) WitResponse
}

type Client struct {
	Token      string
	BaseUrl    string
	ApiVersion string
}

// Convenience function to return an API client with default
// base URL and API version declared.
func NewClient(token string) *Client {
	return &Client{
		Token:      token,
		BaseUrl:    BaseUrl,
		ApiVersion: ApiVersion,
	}
}

// Sends a message to Wit.Ai for parsing. Wit.Ai helps parse
// context (ex. What does a user want?) out of a message
func (c *Client) ParseMessage(message string) WitResponse {
	var witResponse WitResponse

	witBaseUrl := fmt.Sprintf("%s/message?v=%s", c.BaseUrl, c.ApiVersion)
	queryString := url.QueryEscape(message)
	queryUrl := fmt.Sprintf("%s&q=%s", witBaseUrl, queryString)

	req, _ := http.NewRequest("GET", queryUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	request := func(attempt uint) error {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &witResponse)
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
		log.Printf("witai: Failed to retreive wit response")
		witResponse = WitResponse{}
	}

	return witResponse
}
