package wit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var (
	// Base URL for Wit.ai
	BaseUrl = "https://api.wit.ai"

	// Wit.ai API Version
	ApiVersion = "20180128"
)

type client struct {
	Token      string
	BaseUrl    string
	ApiVersion string
}

// Convenience function to return an API client with default
// base URL and API version declared.
func Client(token string) *client {
	return &client{
		Token:      token,
		BaseUrl:    BaseUrl,
		ApiVersion: ApiVersion,
	}
}

// Sends a message to Wit.Ai for parsing. Wit.Ai helps parse
// context (ex. What does a user want?) out of a message
func (c client) ParseMessage(message string) WitResponse {
	witBaseUrl := fmt.Sprintf("%s/message?v=%s", c.BaseUrl, c.ApiVersion)
	queryString := url.QueryEscape(message)
	queryUrl := fmt.Sprintf("%s&q=%s", witBaseUrl, queryString)

	log.Printf("%s", queryUrl)
	req, _ := http.NewRequest("GET", queryUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var witResponse WitResponse

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("witai: response - %s", body)
	jsonErr := json.Unmarshal(body, &witResponse)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return witResponse
}
