package alphapoint

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
)

var BaseUrl = "https://www.alphavantage.co/query"

type Alphapoint interface {
	Convert(fromISO string, toISO string, total float64) (float64, *Conversion)
}

type CurrencyDetails struct {
	Details struct {
		ExchangeRate string `json:"5. Exchange Rate"`
	} `json:"Realtime Currency Exchange Rate"`
}

type Client struct {
	Token   string
	BaseUrl string
}

type Conversion struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
}

// Sets up client to run with AlphaPoint token
func NewClient(token string) *Client {
	return &Client{
		Token:   token,
		BaseUrl: BaseUrl,
	}
}

// Converts from one currency to another using AlphaPoints' API
func (c *Client) Convert(fromISO string, toISO string, total float64) (float64, *Conversion) {
	var currencyDetails CurrencyDetails

	currencyBase := fmt.Sprintf(
		"%s?function=CURRENCY_EXCHANGE_RATE&from_currency=%s&to_currency=%s",
		c.BaseUrl,
		fromISO,
		toISO,
	)
	url := fmt.Sprintf("%s&apikey=%s", currencyBase, c.Token)
	request := func(attempt uint) error {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("alphapoint: response - %s", body)
		jsonErr := json.Unmarshal(body, &currencyDetails)
		if jsonErr != nil {
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
		log.Panicf("alphapoint: unable to convert currency - %s", err)
	}

	exchangeRate, _ := strconv.ParseFloat(currencyDetails.Details.ExchangeRate, 64)
	convertedValue := exchangeRate * total
	conversion := &Conversion{
		fromISO,
		toISO,
		exchangeRate,
	}

	return convertedValue, conversion
}
