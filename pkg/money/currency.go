package money

import (
	// TODO remove this dependency
	m "github.com/Rhymond/go-money"
)

const USD = "USD"

type Currency m.Currency

var defaultCurrency = Currency{
	Code:        "",
	NumericCode: "",
	Fraction:    0,
	Grapheme:    "",
	Template:    "$1",
	Decimal:     "",
	Thousand:    "",
}

func getCurrencyWithDefault(currencyCode string) *Currency {
	currency := m.GetCurrency(currencyCode)

	if currency == nil {
		return &defaultCurrency
	}

	return (*Currency)(currency)
}

func (c *Currency) Equals(oc *Currency) bool {
	return c.Code == oc.Code
}

func (c *Currency) Formatter() *m.Formatter {
	if c == nil {
		return defaultCurrency.Formatter()
	}

	return (*m.Currency)(c).Formatter()
}
