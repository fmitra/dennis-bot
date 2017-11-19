package main

import (
	"testing"
	"time"

	"github.com/kierdavis/dateparser"
)

func TestIsISO(t *testing.T) {
	var isoTests = []struct {
		input string
		expected bool
	}{
		{"USD", true},
		{"PHP", true},
		{"JPY", true},
		{"BTC", true},
		{"SGD", true},
		{"ETH", true},
		{"AAA", false},
	}

	for _, test := range isoTests {
		result := isISO(test.input)
		if result != test.expected {
			t.Errorf("[Input]", test.input, "[Output]", result, "[Expected]", test.expected)
		}
	}
}

func TestParseAmount(t * testing.T) {
	type output struct {
		amount float64
		currency string
	}
	var parseAmountTests = []struct {
		input string
		expected output
	}{
		{"PHP", output{0, "PHP"}},
		{"12SGD", output{12, "SGD"}},
		{"0.00345BTC", output{0.00345, "BTC"}},
		{"200", output{200, ""}},
		{"Hogwarts", output{0, ""}},
	}

	for _, test := range parseAmountTests {
		amount, currency := parseAmount(test.input)
		if amount != test.expected.amount || currency != test.expected.currency {
			t.Error(
				"[Input]", test.input, "[Output]", amount, currency, "[Expected]",
				test.expected.amount, test.expected.currency,
			)
		}
	}
}

func ParseDateTest(t *testing.T) {
	parser := &dateparser.Parser{}
	date, _ := parser.Parse("2017/11/10")
	var parseDateTests = []struct {
		input string
		expected time.Time
	}{
		{"I paid 0.0023BTC on 2017-11-10", date},
		{"2017/11/10 is todays Date", date},
		{"11/10/2017", date},
	}

	for _, test := range parseDateTests {
		result := parseDate(test.input)
		if result != test.expected {
			t.Error(
				"[Input]", test.input, "[Output]", result, "[Expected]", test.expected,
			)
		}
	}
}

func ParseDescriptionTest(t *testing.T) {
	var parseDescriptionTests = []struct {
		input string
		expected string
	}{
		{"for bananas", "bananas"},
		{"apple oranges grapes", "apple oranges grape"},
	}

	for _, test := range parseDescriptionTests {
		result := parseDescription(test.input)
		if result != test.expected {
			t.Error(
				"[Input]", test.input, "[Output]", result, "[Expected]", test.expected,
			)
		}
	}
}
