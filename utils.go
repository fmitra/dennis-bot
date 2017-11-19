package main

import (
	"strings"
	"strconv"
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

// Checks if an input string is a cominbation of an amount and currency
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
