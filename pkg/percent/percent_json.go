package percent

import (
	"encoding/json"
	"errors"
)

var (
	ErrInvalidJSONUnmarshal = errors.New("invalid json unmarshal")
)

// UnmarshalJSON is implementation of json.Unmarshaller
func (p *Percent) UnmarshalJSON(b []byte) error {
	var dataRaw interface{}

	err := json.Unmarshal(b, &dataRaw)

	if err != nil {
		return err
	}

	percentStr, ok := dataRaw.(string)

	if ok {
		pct, err := Parse(percentStr)

		if err != nil {
			return err
		}

		*p = pct

	} else {
		percentFloat, ok := dataRaw.(float64)

		if !ok {
			return ErrInvalidJSONUnmarshal
		}

		*p = Percent(int64(percentFloat * Scale))
	}

	return nil
}

// MarshalText is implementation of encoding.TextMarshaller
// This is needed to correctly encode a map which keys are Percents
func (p Percent) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}

// MarshalJSON is implementation of json.Marshaller
func (p Percent) MarshalJSON() ([]byte, error) {
	return []byte(`"` + p.String() + `"`), nil
}
