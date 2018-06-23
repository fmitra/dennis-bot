package utils

import (
	"testing"
	"time"

	"github.com/kierdavis/dateparser"
	"github.com/stretchr/testify/assert"
)

func TestParseISO(t *testing.T) {
	var isoTests = []struct {
		input    string
		expected string
	}{
		{"USD", "USD"},
		{"php", "PHP"},
		{"JPY", "JPY"},
		{"BTC", "BTC"},
		{" SGD", "SGD"},
		{"ETH", "ETH"},
		{"AAA", ""},
	}

	for _, test := range isoTests {
		result, _ := ParseISO(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestParseAmount(t *testing.T) {
	type output struct {
		amount   float64
		currency string
	}
	var parseAmountTests = []struct {
		input    string
		expected output
	}{
		{"PHP", output{0, "PHP"}},
		{"12SGD", output{12, "SGD"}},
		{"12 SGD", output{12, "SGD"}},
		{"0.00345BTC", output{0.00345, "BTC"}},
		{"200", output{200, ""}},
		{"Hogwarts", output{0, ""}},
	}

	for _, test := range parseAmountTests {
		amount, currency := ParseAmount(test.input)
		assert.Equal(t, test.expected.amount, amount)
		assert.Equal(t, test.expected.currency, currency)
	}
}

func TestDateParser(t *testing.T) {
	parser := &dateparser.Parser{}
	date, _ := parser.Parse("2017/11/10")
	var parseDateTests = []struct {
		input    string
		expected time.Time
	}{
		{"I paid 0.0023BTC on 2017-11-10", date},
		{"2017/11/10 is todays Date", date},
		{"11/10/2017", date},
	}

	for _, test := range parseDateTests {
		result := ParseDate(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestParseDescription(t *testing.T) {
	var parseDescriptionTests = []struct {
		input    string
		expected string
	}{
		{"for bananas", "bananas"},
		{"apple oranges grapes", "apple oranges grapes"},
	}

	for _, test := range parseDescriptionTests {
		result := ParseDescription(test.input)
		assert.Equal(t, test.expected, result)
	}
}
