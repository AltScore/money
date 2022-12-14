package money

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_money_marshal(t *testing.T) {

	amount := FromFloat64(12.34, "MXN")

	s, _ := json.Marshal(amount)

	assert.Equal(t, `{"amount":"12.34","currency":"MXN","display":"$12.34"}`, string(s))
}

func Test_money_marshal_bg_number_shows_no_separators(t *testing.T) {

	amount := FromFloat64(123456.78, "MXN")

	s, _ := json.Marshal(amount)

	assert.Equal(t, `{"amount":"123456.78","currency":"MXN","display":"$123,456.78"}`, string(s))
}

func Test_money_unmarshal(t *testing.T) {
	moneyJson := `{"amount": "128.45", "currency": "ARS", "display":"$123,456.78" }`

	var m Money

	_ = json.Unmarshal([]byte(moneyJson), &m)

	assert.Equal(t, "$128,45", m.String())
}

type Invoice struct {
	TotalAmount Money `json:"totalAmount"`
}

func Test_money_marshal_struct(t *testing.T) {

	invoice := Invoice{
		TotalAmount: NewFromInt(1234, "MXN"),
	}

	s, _ := json.Marshal(invoice)

	fmt.Printf("%s\n", s)

	assert.Equal(t, `{"totalAmount":{"amount":"1234.00","currency":"MXN","display":"$1,234.00"}}`, string(s))
}

func Test_money_unmarshal_struct(t *testing.T) {
	moneyJson := `{"totalAmount":{"amount":"1234.00","currency":"MXN"}}`

	var s Invoice

	_ = json.Unmarshal([]byte(moneyJson), &s)

	assert.Equal(t, "$1,234.00", s.TotalAmount.String())
}

func Test_money_unmarshal_struct_with_money_as_int_number(t *testing.T) {
	moneyJson := `{"totalAmount":{"amount":1234,"currency":"MXN"}}`

	var s Invoice

	_ = json.Unmarshal([]byte(moneyJson), &s)

	assert.Equal(t, "$1,234.00", s.TotalAmount.String())
}

func Test_money_unmarshal_struct_with_money_as_decimal_number(t *testing.T) {
	moneyJson := `{"totalAmount":{"amount":1234.42,"currency":"MXN"}}`

	var s Invoice

	_ = json.Unmarshal([]byte(moneyJson), &s)

	assert.Equal(t, "$1,234.42", s.TotalAmount.String())
}

type Invoice2 struct {
	TotalAmount Money  `json:"totalAmount"`
	Interest    *Money `json:"interest"`
}

func Test_money_marshall_money_and_pointers_to_nil(t *testing.T) {
	invoice := Invoice2{
		TotalAmount: MustParse("42.24", "MXN"),
		Interest:    nil,
	}

	s, err := json.Marshal(invoice)

	assert.Nil(t, err)
	assert.Equal(t, `{"totalAmount":{"amount":"42.24","currency":"MXN","display":"$42.24"},"interest":null}`, string(s))
}

func Test_money_marshall_money_and_pointers_to_money(t *testing.T) {
	interest := MustParse("7.45", "MXN")
	invoice := Invoice2{
		TotalAmount: MustParse("42.24", "MXN"),
		Interest:    &interest,
	}

	s, err := json.Marshal(invoice)

	assert.Nil(t, err)
	assert.Equal(t, `{"totalAmount":{"amount":"42.24","currency":"MXN","display":"$42.24"},"interest":{"amount":"7.45","currency":"MXN","display":"$7.45"}}`, string(s))
}

func Test_money_unmarshall_money_and_pointers_to_nil(t *testing.T) {
	s := `{"totalAmount":{"amount":"42.24","currency":"MXN"},"interest":null}`

	var invoice Invoice2

	err := json.Unmarshal([]byte(s), &invoice)

	assert.Nil(t, err)
	assert.Nil(t, invoice.Interest)
	assert.Equal(t, "$42.24", invoice.TotalAmount.String())
}

func Test_money_unmarshall_money_and_pointers_to_money(t *testing.T) {
	s := `{"totalAmount":{"amount":"42.24","currency":"MXN"},"interest":{"amount":"7.45","currency":"MXN"}}`

	var invoice Invoice2

	err := json.Unmarshal([]byte(s), &invoice)

	assert.Nil(t, err)
	assert.Equal(t, "$42.24", invoice.TotalAmount.String())
	assert.Equal(t, "$7.45", invoice.Interest.String())
}
