package wit

import (
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

var witAiBaseUrl = "https://api.wit.ai"
var witVersion = "20180128"
var Client WitAi

type WitAi struct {
	Token string
}


// Set up client to run with Wit.ai token
func Init(token string) {
	Client = WitAi{token}
}

func (w WitAi) ParseMessage(message string) (WitResponse) {
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
