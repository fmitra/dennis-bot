package alphapoint

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

var baseUrl = "https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&"

var Client AlphaPoint

type AlphaPoint struct {
	Token string
}

type CurrencyDetails struct {
	Details struct {
		ExchangeRate string `json:"5. Exchange Rate"`
	} `json:"Realtime Currency Exchange Rate"`
}

// Sets up client to run with AlphaPoint token
func Init(token string) {
	Client = AlphaPoint{token}
}

// Converts from one currency to another using AlphaPoints' API
func (a AlphaPoint) Convert(fromISO string, toISO string, total float64) (float64) {
	currencyBase := fmt.Sprintf(
		"%sfrom_currency=%s&to_currency=%s", baseUrl, fromISO, toISO,
	)
	url := fmt.Sprintf("%s&apikey=%s", currencyBase, a.Token)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	var currencyDetails CurrencyDetails
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("alphapoint: response - %s", body)
	jsonErr := json.Unmarshal(body, &currencyDetails)
	if jsonErr != nil {
		panic(err)
	}

	exchangeRate, _ := strconv.ParseFloat(currencyDetails.Details.ExchangeRate, 64)
	conversion := exchangeRate * total
	return conversion
}
