package currency

import (
	// TODO remove this dependency

	"strconv"
	"strings"
)

// Currency represents money currency information required for formatting.
type Currency struct {
	Code        string
	NumericCode string
	Fraction    int
	Grapheme    string
	Template    string
	Decimal     string
	Thousand    string
}

// Get returns the currency given the code.
func Get(code string) *Currency {
	return currencies.CurrencyByCode(code)
}

// GetOrDefault returns the currency given the code or default currency if not found.
func GetOrDefault(currencyCode string) *Currency {
	code := strings.ToUpper(currencyCode)

	currency := Get(code)

	if currency == nil {
		currency = &Currency{
			Code:     code,
			Template: "$1",
			Grapheme: code,
			Decimal:  ".",
			Thousand: ",",
			Fraction: 0,
		}
		currencies.Add(currency) // Cache the currency
		return currency
	}

	return currency
}

func (c *Currency) Equals(oc *Currency) bool {
	if c == nil || oc == nil {
		return c == oc
	}
	return c.Code == oc.Code
}

func (c *Currency) Format(amount int64) string { // TODO improve, use better formatter
	if c == nil {
		return GetOrDefault("").Format(amount)
	}
	// Work with absolute amount value
	sa := strconv.FormatInt(abs(amount), 10)

	if len(sa) <= c.Fraction {
		sa = strings.Repeat("0", c.Fraction-len(sa)+1) + sa
	}

	if c.Thousand != "" {
		for i := len(sa) - c.Fraction - 3; i > 0; i -= 3 {
			sa = sa[:i] + c.Thousand + sa[i:]
		}
	}

	if c.Fraction > 0 {
		sa = sa[:len(sa)-c.Fraction] + c.Decimal + sa[len(sa)-c.Fraction:]
	}
	sa = strings.Replace(c.Template, "1", sa, 1)
	sa = strings.Replace(sa, "$", c.Grapheme, 1)

	// Add minus sign for negative amount.
	if amount < 0 {
		sa = "-" + sa
	}

	return sa
}

func abs(amount int64) int64 {
	if amount < 0 {
		return -amount
	}

	return amount
}
