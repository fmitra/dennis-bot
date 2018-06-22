// Package alphapoint implements a wrapper for the Alphapoint API.
// AlphaPoint provides an API for currency conversion.
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

// BaseURL for Alphapoint
var BaseURL = "https://www.alphavantage.co/query"

// Alphapoint is an itnerface to provide utility methods to interact
// with the Alphapoint API.
type Alphapoint interface {
	Convert(fromISO string, toISO string, total float64) (float64, *Conversion)
}

// CurrencyDetails describes the exchange rate of a currency.
type CurrencyDetails struct {
	Details struct {
		ExchangeRate string `json:"5. Exchange Rate"`
	} `json:"Realtime Currency Exchange Rate"`
}

// Client is a consumer of the Alphapoint API.
type Client struct {
	Token   string
	BaseURL string
}

// Conversion describes the exchange rate between two currencies.
type Conversion struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
}

// NewClient returns a Client with default BaseUrl.
func NewClient(token string) *Client {
	return &Client{
		Token:   token,
		BaseURL: BaseURL,
	}
}

// Convert converts from one currency to another using AlphaPoints' API.
func (c *Client) Convert(fromISO string, toISO string, total float64) (float64, *Conversion) {
	var currencyDetails CurrencyDetails

	currencyBase := fmt.Sprintf(
		"%s?function=CURRENCY_EXCHANGE_RATE&from_currency=%s&to_currency=%s",
		c.BaseURL,
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
		err = json.Unmarshal(body, &currencyDetails)
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
