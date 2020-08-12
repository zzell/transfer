package model

import "errors"

// Currency enum
type Currency int

// available currencies
const (
	BTC Currency = iota
	ETH
)

var currencies = [...]string{"BTC", "ETH"}

// String converts currency enum to string value
func (c Currency) String() string {
	return currencies[c]
}

// CurrencyFromString converts string to enum type
func CurrencyFromString(s string) (Currency, error) {
	for i, c := range currencies {
		if s == c {
			return Currency(i), nil
		}
	}

	return 0, errors.New("illegal currency type")
}
