package main

import (
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
)

var witAiBaseUrl = "https://api.wit.ai"
var witVersion = "20180128"
var witAi = WitAi{os.Getenv("WITAI_AUTH_TOKEN")}

type WitAi struct {
	Token string
}

func (w WitAi) parseMessage(message string) (WitResponse) {
	baseUrl := fmt.Sprintf("%s/message?v=%s", witAiBaseUrl, witVersion)
	queryString := url.QueryEscape(message)
	queryUrl := fmt.Sprintf("%s&q=%s", baseUrl, queryString)

	log.Printf("%s", queryUrl)
	req, _ := http.NewRequest("GET", queryUrl, nil)
	req.Header.Set("Authorization", w.Token)

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
		panic(err)
	}

	return witResponse
}
