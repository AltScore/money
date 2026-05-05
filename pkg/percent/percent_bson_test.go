package percent

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type sample struct {
	Percent Percent `bson:"directPercent"`
}

func Test_MarshalBSON_is_the_inverse_of_UnmarshallBSON(t *testing.T) {
	reg := newTestRegistry()

	values := []string{
		"0",
		"10",
		"0.1",
		"12345",
		"123.678",
		"0.045",
		"-10.42",
	}

	for _, str := range values {
		value := MustParse(str)

		original := sample{value}

		var buf bytes.Buffer
		vw := bson.NewDocumentWriter(&buf)
		enc := bson.NewEncoder(vw)
		enc.SetRegistry(reg)
		err := enc.Encode(&original)
		require.Nil(t, err)

		vr := bson.NewDocumentReader(bytes.NewReader(buf.Bytes()))
		dec := bson.NewDecoder(vr)
		dec.SetRegistry(reg)

		var decoded sample
		err = dec.Decode(&decoded)

		require.Nil(t, err)
		require.Equal(t, value, decoded.Percent)
		require.Equal(t, str, decoded.Percent.String())
	}
}

func newTestRegistry() *bson.Registry {
	reg := bson.NewRegistry()
	RegisterPercentBSONCodec(reg)
	return reg
}