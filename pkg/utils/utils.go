// Package utils provide common utility methods to share across packages.
package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/kierdavis/dateparser"
)

// ParseISO checks if an input string is a valid currency ISO.
func ParseISO(s string) (string, error) {
	cleanString := strings.Replace(s, " ", "", -1)
	upperS := strings.ToUpper(cleanString)
	for _, iso := range CURRENCIES {
		if upperS == iso {
			return upperS, nil
		}
	}
	return "", errors.New("invalid currency")
}

// ParseAmount checks a string for possible currency and amount values,
// for example "100 USD" should return 100 as a float and USD as
// the currency ISO.
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

// ParseDate checks if a string is a possible date value.
func ParseDate(s string) (inferredDate time.Time) {
	lowerCase := strings.ToLower(s)
	splitString := strings.Split(lowerCase, " ")
	today := time.Now()
	date := today

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

// ParseDescription checks if a string is a description.
func ParseDescription(s string) (description string) {
	lowerCase := strings.ToLower(s)
	if strings.HasPrefix(lowerCase, "for ") {
		return s[4:]
	}

	return s
}
