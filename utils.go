package main

import (
	"strings"
	"strconv"
	"time"

	"github.com/kierdavis/dateparser"
)

// Checks if an input string is a currency ISO
func isISO(s string) (bool) {
	upperS := strings.ToUpper(s)
	for _, iso := range CURRENCIES {
		if upperS == iso {
			return true
		}
	}
	return false
}

// Checks if an input string contains a cominbation of
// an amount and currency. For simplicity, we require
// an amount to be written on the left of the currency ISO
func inferAmount(s string) (float64, string) {
	cleanString := strings.Replace(s, " ", "", -1)
	upperCase := strings.ToUpper(cleanString)

	if isISO(cleanString) {
		return 0, upperCase
	}

	var currency string
	for _, iso := range CURRENCIES {
		if strings.Contains(upperCase, iso) {
			currency = iso

			// We expect the amount to be written on the left of the currency
			splitString := strings.Split(upperCase, iso)
			amount := strings.Replace(splitString[0], ",", "", -1)
			parsedAmount, err := strconv.ParseFloat(amount, 64)

			if err != nil {
				return 0, currency
			}
			return parsedAmount, currency
		}
	}
	return 0, ""
}

// Checks if an input string contains a possible date value
func inferDate(s string) (time.Time) {
	lowerCase := strings.ToLower(s)
	splitString := strings.Split(lowerCase, " ")
	today := time.Now()
	var date time.Time

	date = today
	if strings.Contains(lowerCase, "yesterday") {
		date = today.AddDate(0, 0, -1)
	}

	// Check if any item in the split is an actual date string
	parser := &dateparser.Parser{}
	for _, item := range splitString {
		parsedTime, err := parser.Parse(item)
		if err == nil {
			date = parsedTime
		}
	}

	return date
}
