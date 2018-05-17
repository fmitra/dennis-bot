package wit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
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

// Set up client to run with Wit.ai token
func Client(token string) *client {
	return &client{
		Token: token,
		BaseUrl: BaseUrl,
		ApiVersion: ApiVersion,
	}
}

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
