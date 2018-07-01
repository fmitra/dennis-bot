// Package wit implements a wrapper for the Wit.ai API.
// Wit provides an NLP API for parsing context from human text.
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

const (
	// BaseURL for Wit.ai
	BaseURL = "https://api.wit.ai"

	// APIVersion is the Wit.ai API Version
	APIVersion = "20180128"
)

// Wit is an itnerface to provide utility methods to interact with
// the Wit.ai API.
type Wit interface {
	ParseMessage(message string) Response
}

// Client is a consumer of the Wit.ai API.
type Client struct {
	Token      string
	BaseURL    string
	APIVersion string
}

// NewClient returns a Client with a default BaseURL.
func NewClient(token string) *Client {
	return &Client{
		Token:      token,
		BaseURL:    BaseURL,
		APIVersion: APIVersion,
	}
}

// ParseMessage passes a message to Wit.ai to infer context, for example,
// Wit.ai may infer that the user is trying to track an expense.
// Any error from Wit.ai is simply handled as an empty Response indicating,
// we are not able to infer anything.
func (c *Client) ParseMessage(message string) Response {
	var response Response

	witBaseURL := fmt.Sprintf("%s/message?v=%s", c.BaseURL, c.APIVersion)
	queryString := url.QueryEscape(message)
	queryURL := fmt.Sprintf("%s&q=%s", witBaseURL, queryString)

	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return response
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	request := func(attempt uint) error {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &response)

		// There is no need to retry if we receive a successful request containing
		// invalid data. We simply treat it as a empty response and carry on.
		if err != nil {
			log.Printf("wit: invalid response received - %s", err)
		}
		return nil
	}

	err = retry.Retry(
		request,
		strategy.Limit(10),
		strategy.Backoff(backoff.Exponential(time.Second, 2)),
	)

	if err != nil {
		log.Printf("witai: failed to retrieve wit response")
	}

	return response
}
