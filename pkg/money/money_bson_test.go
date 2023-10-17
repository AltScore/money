package money

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_MarshalBSON_is_the_inverse_of_UnmarshallBSON(t *testing.T) {
	values := []Money{
		NewFromInt(0, "MXN"),
		NewFromInt(1000, "USD"),
		NewFromInt(1500, "ARS"), // Have comma as decimal separator
	}

	for _, value := range values {

		bytes, err := bson.Marshal(value)

		require.Nil(t, err)

		var decoded Money
		err = bson.Unmarshal(bytes, &decoded)

		require.Nil(t, err)
		require.Equal(t, value, decoded)
	}
}

type sample struct {
	DirectMoney  Money  `bson:"directMoney"`
	PointerMoney *Money `bson:"pointerMoney"`
}

func Test_MarshalBSON_is_the_inverse_of_UnmarshallBSON_with_structs(t *testing.T) {
	pointerToMoney := NewFromInt(1000, "USD")
	values := []sample{
		{
			NewFromInt(3243, "MXN"),
			&pointerToMoney,
		},
	}

	for _, value := range values {

		bytes, err := bson.Marshal(&value)

		require.Nil(t, err)

		var m bson.D

		_ = bson.Unmarshal(bytes, &m)

		fmt.Println(m)

		var decoded sample
		err = bson.Unmarshal(bytes, &decoded)

		require.Nil(t, err)
		require.Equal(t, value, decoded)
	}
}

func TestMoney_MarshalBSON(t *testing.T) {
	tests := []struct {
		name    string
		value   Money
		want    any
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:  "Zero",
			value: Zero("USD"),
			want: map[string]interface{}{
				"amount":   "0.00",
				"currency": "USD",
			},
			wantErr: assert.NoError,
		},
		{
			name:  "Positive",
			value: FromFloat64(1000.42, "ARS"),
			want: map[string]interface{}{
				"amount":   "1000.42",
				"currency": "ARS",
			},
			wantErr: assert.NoError,
		},
		{
			name:  "Negative",
			value: FromFloat64(-328.76, "MXN"),
			want: map[string]interface{}{
				"amount":   "-328.76",
				"currency": "MXN",
			},
			wantErr: assert.NoError,
		},
		{
			name:  "Small Negative cents",
			value: FromFloat64(-0.01, "MXN"),
			want: map[string]interface{}{
				"amount":   "-0.01",
				"currency": "MXN",
			},
			wantErr: assert.NoError,
		},
		{
			name:  "Small Negative tens",
			value: FromFloat64(-0.10, "MXN"),
			want: map[string]interface{}{
				"amount":   "-0.10",
				"currency": "MXN",
			},
			wantErr: assert.NoError,
		},
		{
			name:  "Small Positive",
			value: FromFloat64(0.01, "MXN"),
			want: map[string]interface{}{
				"amount":   "0.01",
				"currency": "MXN",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := tt.value.MarshalBSON()

			if !tt.wantErr(t, err, fmt.Sprintf("MarshalBSON()")) {
				return
			}

			var got map[string]interface{}

			err = bson.Unmarshal(bytes, &got)

			require.NoError(t, err)

			assert.Equalf(t, tt.want, got, "MarshalBSON()")
		})
	}
}
