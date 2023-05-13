package currency

import (
	"sync"
)

type Currencies map[string]*Currency

// currenciesLock is used to forbid concurrent modification to currencies.
var currenciesLock sync.RWMutex

// CurrencyByNumericCode returns the currency given the numeric code defined in ISO-4271.
func (c Currencies) CurrencyByNumericCode(code string) *Currency {
	currenciesLock.RLock()
	defer currenciesLock.RUnlock()

	for _, sc := range c {
		if sc.NumericCode == code {
			return sc
		}
	}

	return nil
}

// CurrencyByCode returns the currency given the currency code defined as a constant.
func (c Currencies) CurrencyByCode(code string) *Currency {
	currenciesLock.RLock()
	defer currenciesLock.RUnlock()

	sc, ok := c[code]

	if !ok {
		return nil
	}

	return sc
}

// Add updates currencies list by adding a given Currency to it.
func (c Currencies) Add(currency *Currency) Currencies {
	currenciesLock.Lock()
	defer currenciesLock.Unlock()

	c[currency.Code] = currency
	return c
}

// AddCurrency lets you insert or update currency in currencies list.
func AddCurrency(code, Grapheme, Template, Decimal, Thousand string, Fraction int) *Currency {
	c := Currency{
		Code:     code,
		Grapheme: Grapheme,
		Template: Template,
		Decimal:  Decimal,
		Thousand: Thousand,
		Fraction: Fraction,
	}
	currencies.Add(&c)
	return &c
}
