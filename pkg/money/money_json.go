package money

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/AltScore/money/pkg/parsers"
	m "github.com/Rhymond/go-money"
	"go.uber.org/zap"
)

// UnmarshalJSON is implementation of json.Unmarshaller
func (a *Money) UnmarshalJSON(b []byte) error {
	data := make(map[string]interface{})
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	currencyCode, err := jsonExtractCurrency(data)

	if err != nil {
		return err
	}

	amount, err := jsonExtractAmount(data, currencyCode)

	if err != nil {
		return err
	}

	var ref *m.Money
	if amount == 0 && currencyCode == "" {
		ref = &m.Money{}
	} else {
		ref = m.New(amount, currencyCode)
	}

	*a = Money(*ref)
	return nil

}

func jsonExtractAmount(data map[string]interface{}, currencyCode string) (int64, error) {
	amountRaw, ok := data["amount"]
	if !ok {
		return 0, nil
	}

	currency := m.GetCurrency(currencyCode)

	if amountStr, ok := amountRaw.(string); ok {

		amount, err := parsers.ParseNumber(amountStr, currency.Fraction)

		if err != nil {
			return 0, m.ErrInvalidJSONUnmarshal
		}

		return amount, nil
	}

	// It is expressed as a number
	amountFloat, ok := amountRaw.(float64)
	if !ok {
		return 0, m.ErrInvalidJSONUnmarshal
	}

	return float2EquivalentInt(amountFloat, currency), nil
}

func jsonExtractCurrency(data map[string]interface{}) (string, error) {
	if currencyRaw, ok := data["currency"]; !ok {
		return "", nil
	} else if currencyCode, ok := currencyRaw.(string); !ok {
		return "", m.ErrInvalidJSONUnmarshal
	} else {
		return currencyCode, nil
	}
}

// MarshalJSON is implementation of json.Marshaller
func (a Money) MarshalJSON() ([]byte, error) {
	ma := a.asMoney()

	var jsonValue string

	if ma.Currency() == nil {
		jsonValue = fmt.Sprintf(`{"amount":"%d","currency":"?","display":"%d"}`, ma.Amount(), ma.Amount())
	} else {
		currencyCode, amountStr := formatAsNumber(ma)

		jsonValue = fmt.Sprintf(`{"amount":"%s","currency":"%s","display":"%s"}`, amountStr, currencyCode, ma.Display())
	}

	return []byte(jsonValue), nil
}

// move this to formatter in github.com/Rhymond/go-money
func formatAsNumber(ma *m.Money) (string, string) {
	currency := ma.Currency()

	if currency == nil {
		if ma.Amount() != 0 {
			amount := strconv.FormatInt(ma.Amount(), 10)
			zap.L().Warn("Currency is nil, amount is " + amount)
		}
		currency = defaultCurrency
	}

	formatter := *currency.Formatter()

	formatter.Grapheme = "" // Remove grapheme
	formatter.Thousand = "" // Remove thousand-separator
	amountStr := formatter.Format(ma.Amount())
	return currency.Code, amountStr
}
