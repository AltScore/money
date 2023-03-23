package percent

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

type sample struct {
	Percent Percent `bson:"directPercent"`
}

func Test_MarshalBSON_is_the_inverse_of_UnmarshallBSON(t *testing.T) {
	registerCodex()

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

		bytes, err := bson.Marshal(&original)

		require.Nil(t, err)

		var decoded sample
		err = bson.Unmarshal(bytes, &decoded)

		require.Nil(t, err)
		require.Equal(t, value, decoded.Percent)
		require.Equal(t, str, decoded.Percent.String())
	}
}

func registerCodex() {
	codec, err := bsoncodec.NewStructCodec(bsoncodec.JSONFallbackStructTagParser)

	if err != nil {
		panic(err)
	}

	builder := bson.NewRegistryBuilder()
	builder.RegisterDefaultEncoder(reflect.Struct, codec)
	builder.RegisterDefaultDecoder(reflect.Struct, codec)

	RegisterPercentBSONCodec(bsonRegistryBuilderAdapter{builder})
	bson.DefaultRegistry = builder.Build()
}

type bsonRegistryBuilderAdapter struct {
	*bsoncodec.RegistryBuilder
}

func (b bsonRegistryBuilderAdapter) RegisterTypeEncoder(t reflect.Type, dec bsoncodec.ValueEncoder) {
	b.RegistryBuilder.RegisterTypeEncoder(t, dec)
}

func (b bsonRegistryBuilderAdapter) RegisterTypeDecoder(t reflect.Type, dec bsoncodec.ValueDecoder) {
	b.RegistryBuilder.RegisterTypeDecoder(t, dec)
}
