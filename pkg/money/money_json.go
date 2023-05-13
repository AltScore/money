package money

import (
	"encoding/json"
	"fmt"

	currency2 "github.com/AltScore/money/pkg/money/currency"
	"github.com/AltScore/money/pkg/parsers"
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

	var ref Money
	if amount == 0 && currencyCode == "" {
		ref = Money{}
	} else {
		ref = fromEquivalentInt(amount, currencyCode)
	}

	*a = ref
	return nil

}

func jsonExtractAmount(data map[string]interface{}, currencyCode string) (int64, error) {
	amountRaw, ok := data["amount"]
	if !ok {
		return 0, nil
	}

	currency := currency2.GetOrDefault(currencyCode)

	if amountStr, ok := amountRaw.(string); ok {

		amount, err := parsers.ParseNumber(amountStr, currency.Fraction)

		if err != nil {
			return 0, ErrInvalidJSONUnmarshal
		}

		return amount, nil
	}

	// It is expressed as a number
	amountFloat, ok := amountRaw.(float64)
	if !ok {
		return 0, ErrInvalidJSONUnmarshal
	}

	return float2EquivalentInt(amountFloat, currency), nil
}

func jsonExtractCurrency(data map[string]interface{}) (string, error) {
	if currencyRaw, ok := data["currency"]; !ok {
		return "", nil
	} else if currencyCode, ok := currencyRaw.(string); !ok {
		return "", ErrInvalidJSONUnmarshal
	} else {
		return currencyCode, nil
	}
}

// MarshalJSON is implementation of json.Marshaller
func (a Money) MarshalJSON() ([]byte, error) {
	var jsonValue string

	if a.currency == nil {
		jsonValue = fmt.Sprintf(`{"amount":"%d","currency":"?","display":"%d"}`, a.amount, a.amount)
	} else {
		currencyCode, amountStr := a.formatAsNumber()

		jsonValue = fmt.Sprintf(`{"amount":"%s","currency":"%s","display":"%s"}`, amountStr, currencyCode, a.String())
	}

	return []byte(jsonValue), nil
}
