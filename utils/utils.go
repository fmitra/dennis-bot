package utils

import (
	"strings"
	"strconv"
	"time"

	"github.com/kierdavis/dateparser"
)

// Checks if an input string is a currency ISO
func isISO(s string) (bool) {
	cleanString := strings.Replace(s, " ", "", -1)
	upperS := strings.ToUpper(cleanString)
	for _, iso := range CURRENCIES {
		if upperS == iso {
			return true
		}
	}
	return false
}

// Parses an input string for an amount and currency
func ParseAmount(s string) (amount float64, currency string) {
	upperCase := strings.ToUpper(s)
	for _, currency := range CURRENCIES {
		if strings.Contains(upperCase, currency) {

			amount := strings.Split(upperCase, currency)[0]
			cleanAmount := strings.Replace(amount, " ", "", -1)
			parsedAmount, err := strconv.ParseFloat(cleanAmount, 64)

			if err != nil {
				return 0, currency
			}

			return parsedAmount, currency
		}
	}

	parsedAmount, _ := strconv.ParseFloat(s, 64)
	return parsedAmount, ""
}

// Parses an input string for a date value
func ParseDate(s string) (inferredDate time.Time) {
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

// Parses an input string for a description
func ParseDescription(s string) (description string) {
	lowerCase := strings.ToLower(s)
	if strings.HasPrefix(lowerCase, "for ") {
		return s[3:]
	}

	return s
}
