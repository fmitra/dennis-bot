package alphapoint

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var Client AlphaPoint

const baseUrl = "https://www.alphavantage.co/query"

type HttpLib interface {
	Get(url string) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type AlphaPoint struct {
	Token   string
	BaseUrl string
	Http    HttpLib
}

type CurrencyDetails struct {
	Details struct {
		ExchangeRate string `json:"5. Exchange Rate"`
	} `json:"Realtime Currency Exchange Rate"`
}

// Sets up client to run with AlphaPoint token
func Init(token string, httpLib HttpLib) {
	Client = AlphaPoint{
		token,
		baseUrl,
		httpLib,
	}
}

// Converts from one currency to another using AlphaPoints' API
func (a AlphaPoint) Convert(fromISO string, toISO string, total float64) float64 {
	currencyBase := fmt.Sprintf(
		"%s?function=CURRENCY_EXCHANGE_RATE&from_currency=%s&to_currency=%s",
		a.BaseUrl,
		fromISO,
		toISO,
	)
	url := fmt.Sprintf("%s&apikey=%s", currencyBase, a.Token)

	resp, err := a.Http.Get(url)
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
