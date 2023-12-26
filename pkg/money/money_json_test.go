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

func Test_money_marshall_in_ARS_shows_point_as_decimal_separator(t *testing.T) {

	amount := FromFloat64(12.34, "ARS")

	s, _ := json.Marshal(amount)

	assert.Equal(t, `{"amount":"12.34","currency":"ARS","display":"$12,34"}`, string(s))
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

func Test_money_unmarshal_struct_using_ARS(t *testing.T) {
	moneyJson := `{"totalAmount":{"amount":"1234.00","currency":"ARS"}}`

	var s Invoice

	_ = json.Unmarshal([]byte(moneyJson), &s)

	assert.Equal(t, "$1.234,00", s.TotalAmount.String())
	assert.Equal(t, "1234.00", s.TotalAmount.Amount())
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

func TestMoney_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		m    Money
		want string
	}{
		{
			name: "empty",
			m:    Money{},
			want: `{"amount":"0","currency":"?","display":"0"}`,
		},
		{
			name: "zero",
			m:    FromFloat64(0, "MXN"),
			want: `{"amount":"0.00","currency":"MXN","display":"$0.00"}`,
		},
		{
			name: "zero in ARS",
			m:    FromFloat64(0, "ARS"),
			want: `{"amount":"0.00","currency":"ARS","display":"$0,00"}`,
		},
		{
			name: "cents",
			m:    FromFloat64(0.01, "MXN"),
			want: `{"amount":"0.01","currency":"MXN","display":"$0.01"}`,
		},
		{
			name: "negative cents",
			m:    FromFloat64(-0.01, "MXN"),
			want: `{"amount":"-0.01","currency":"MXN","display":"-$0.01"}`,
		},
		{
			name: "cents in ARS",
			m:    FromFloat64(0.04, "ARS"),
			want: `{"amount":"0.04","currency":"ARS","display":"$0,04"}`,
		},
		{
			name: "negative cents in ARS",
			m:    FromFloat64(-0.04, "ARS"),
			want: `{"amount":"-0.04","currency":"ARS","display":"-$0,04"}`,
		},
		{
			name: "dimes",
			m:    FromFloat64(0.23, "MXN"),
			want: `{"amount":"0.23","currency":"MXN","display":"$0.23"}`,
		},
		{
			name: "dimes in ARS",
			m:    FromFloat64(0.42, "ARS"),
			want: `{"amount":"0.42","currency":"ARS","display":"$0,42"}`,
		},
		{
			name: "with integral part",
			m:    FromFloat64(123.23, "MXN"),
			want: `{"amount":"123.23","currency":"MXN","display":"$123.23"}`,
		},
		{
			name: "with integral part in ARS",
			m:    FromFloat64(718.64, "ARS"),
			want: `{"amount":"718.64","currency":"ARS","display":"$718,64"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.MarshalJSON()
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, string(got), "MarshalJSON()")
		})
	}
}

func TestMoney_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    Money
		wantErr bool
	}{
		{
			name:    "empty",
			data:    `{}`,
			wantErr: true,
		},
		{
			name:    "missing amount",
			data:    `{"currency":"MXN"}`,
			wantErr: true,
		},
		{
			name:    "missing currency",
			data:    `{"amount":"123.45"}`,
			wantErr: true,
		},
		{
			name:    "invalid amount two decimal places",
			data:    `{"amount":"123.45.67","currency":"MXN"}`,
			wantErr: true,
		},
		{
			name:    "invalid amount decimal separator",
			data:    `{"amount":"123,45","currency":"MXN"}`,
			wantErr: true,
		},
		{
			name:    "invalid amount symbol",
			data:    `{"amount":"$123.45","currency":"MXN"}`,
			wantErr: true,
		},
		{
			name:    "invalid currency",
			data:    `{"amount":"123.45","currency":"MNX"}`,
			wantErr: true,
		},
		{
			name:    "invalid currency type",
			data:    `{"amount":"123.45","currency":123}`,
			wantErr: true,
		},
		{
			name: "zero",
			data: `{"amount":"0.00","currency":"MXN"}`,
			want: FromFloat64(0, "MXN"),
		},
		{
			name: "zero in ARS",
			data: `{"amount":"0.00","currency":"ARS"}`,
			want: FromFloat64(0, "ARS"),
		},
		{
			name: "cents",
			data: `{"amount":"0.01","currency":"MXN"}`,
			want: FromFloat64(0.01, "MXN"),
		},
		{
			name: "cents in ARS",
			data: `{"amount":"0.04","currency":"ARS"}`,
			want: FromFloat64(0.04, "ARS"),
		},
		{
			name:    "negative amount",
			data:    `{"amount":"-123.45","currency":"MXN"}`,
			wantErr: true,
		},
		{
			name:    "negative cents",
			data:    `{"amount":"-0.01","currency":"MXN"}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Money
			if err := got.UnmarshalJSON([]byte(tt.data)); (err != nil) != tt.wantErr {
				t.Errorf("Money.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
