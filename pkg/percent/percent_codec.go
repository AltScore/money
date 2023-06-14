package percent

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type BsonRegistryBuilder interface {
	RegisterTypeEncoder(t reflect.Type, dec bsoncodec.ValueEncoder)
	RegisterTypeDecoder(t reflect.Type, dec bsoncodec.ValueDecoder)
}

// Codec is the Codec used for percent.Percent values.
type Codec struct {
	typeOf reflect.Type
}

var (
	emptyValue = reflect.Value{}

	defaultPercentCodec = NewPercentCodec()

	_ bsoncodec.ValueCodec   = defaultPercentCodec
	_ bsoncodec.ValueDecoder = defaultPercentCodec
)

// RegisterPercentBSONCodec register in the BSON registry a Codec to handle objects values of type percent.Percent
// Prefer to use Register method
func RegisterPercentBSONCodec(builder BsonRegistryBuilder) {
	codec := NewPercentCodec()

	builder.RegisterTypeEncoder(codec.typeOf, codec)
	builder.RegisterTypeDecoder(codec.typeOf, codec)
}

// NewPercentCodec returns a PercentCodec with options opts.
func NewPercentCodec() *Codec {
	return &Codec{
		typeOf: reflect.TypeOf(Zero),
	}
}

func (pc *Codec) Register(registryBuilder *bsoncodec.RegistryBuilder) {
	registryBuilder.RegisterTypeEncoder(pc.typeOf, pc)
	registryBuilder.RegisterTypeDecoder(pc.typeOf, pc)
}

//nolint:cyclop // this is a simple switch for type matching
func (pc *Codec) decodeType(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) {
	if t != pc.typeOf {
		return emptyValue, bsoncodec.ValueDecoderError{
			Name:     "PercentDecodeValue",
			Types:    []reflect.Type{pc.typeOf},
			Received: reflect.Zero(t),
		}
	}

	var percentVal Percent
	switch vrType := vr.Type(); vrType {
	case bsontype.String:
		// assume strings are in the "9999.9999" format
		percentStr, err := vr.ReadString()
		if err != nil {
			return emptyValue, err
		}
		percentVal, err = Parse(percentStr)
		if err != nil {
			return emptyValue, err
		}
	case bsontype.Int64:
		i64, err := vr.ReadInt64()
		if err != nil {
			return emptyValue, err
		}
		percentVal = Percent(i64)
	case bsontype.Int32:
		i32, err := vr.ReadInt32()
		if err != nil {
			return emptyValue, err
		}
		percentVal = Percent(i32)
	case bsontype.Null:
		if err := vr.ReadNull(); err != nil {
			return emptyValue, err
		}
	case bsontype.Undefined:
		if err := vr.ReadUndefined(); err != nil {
			return emptyValue, err
		}
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a percent.Percent", vrType)
	}

	return reflect.ValueOf(percentVal), nil
}

// DecodeValue is the ValueDecoderFunc for time.Time.
func (pc *Codec) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != pc.typeOf {
		return bsoncodec.ValueDecoderError{Name: "PercentDecodeValue", Types: []reflect.Type{pc.typeOf}, Received: val}
	}

	elem, err := pc.decodeType(dc, vr, pc.typeOf)
	if err != nil {
		return err
	}

	val.Set(elem)
	return nil
}

// EncodeValue is the ValueEncoderFunc for time.TIme.
func (pc *Codec) EncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != pc.typeOf {
		return bsoncodec.ValueEncoderError{Name: "PercentEncodeValue", Types: []reflect.Type{pc.typeOf}, Received: val}
	}
	p := val.Interface().(Percent) //nolint:forcetypeassert // previous check ensures this is a Percent
	return vw.WriteString(p.String())
}
