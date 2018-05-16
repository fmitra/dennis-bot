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
	baseUrl = "https://api.wit.ai"

	// Wit.ai API Version
	apiVersion = "20180128"
)

var Client WitAi

type WitAi struct {
	Token      string
	BaseUrl    string
	ApiVersion string
}

// Set up client to run with Wit.ai token
func Init(token string) {
	Client = WitAi{
		token,
		baseUrl,
		apiVersion,
	}
}

func (w WitAi) ParseMessage(message string) WitResponse {
	witBaseUrl := fmt.Sprintf("%s/message?v=%s", w.BaseUrl, w.ApiVersion)
	queryString := url.QueryEscape(message)
	queryUrl := fmt.Sprintf("%s&q=%s", witBaseUrl, queryString)

	log.Printf("%s", queryUrl)
	req, _ := http.NewRequest("GET", queryUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.Token))

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
