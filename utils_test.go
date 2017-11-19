package main

import "testing"

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
			t.Errorf("isISO(%d) = %d, want %d", test.input, result, test.expected)
		}
	}
}

func TestInferAmount(t *testing.T) {
	type output struct {
		amount float64
		currency string
	}
	var inferAmountTests = []struct {
		input string
		expected output
	}{
		{"PHP", output{0, "PHP"}},
		{"   SGD", output{0, "SGD"}},
		{"12 SGD", output{12, "SGD"}},
		{"10,000.50 JPY", output{10000.50, "JPY"}},
		{"0.000345BTC", output{0.000345, "BTC"}},
		{"12USD", output{12, "USD"}},
		{"12,000USD Something here", output{12000, "USD"}},
		{"0.0023BTC on 2017/12/18", output{0.0023, "BTC"}},
		{"I paid 0.0023BTC on 2017/12/18", output{0, "BTC"}},
	}

	for _, test := range inferAmountTests {
		amount, currency := inferAmount(test.input)
		e_amount := test.expected.amount
		e_currency := test.expected.currency
		if amount != e_amount || currency != e_currency {
			t.Error("inferAmount = %d, %d, want %d, %d", amount, currency, e_amount, e_currency)
		}
	}
}
