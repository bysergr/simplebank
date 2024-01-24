package utils

import "slices"

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

var validCurrencies = []string{USD, EUR, CAD}

func IsValidCurrency(currency string) bool {
	if slices.Contains(validCurrencies, currency) {
		return true
	}

	return false
}
