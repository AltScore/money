package money

import (
	"fmt"
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
