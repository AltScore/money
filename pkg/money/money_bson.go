package money

import (
	"errors"
	"strings"

	"github.com/AltScore/money/pkg/money/currency"
	"github.com/AltScore/money/pkg/parsers"
	"go.mongodb.org/mongo-driver/bson"
)

var ErrInvalidBSONUnmarshal = errors.New("invalid bson unmarshal")

type bsonMoney struct {
	Amount   string `bson:"amount"`
	Currency string `bson:"currency"`
}

// UnmarshalBSON is implementation of json.Unmarshaller
func (a *Money) UnmarshalBSON(b []byte) error {
	bm := bsonMoney{}

	err := bson.Unmarshal(b, &bm)

	if err != nil {
		return err
	}

	cur := currency.GetOrDefault(bm.Currency)

	am := bm.Amount

	if cur.Decimal != "." {
		am = strings.ReplaceAll(am, cur.Decimal, ".")
	}

	amount, err := parsers.ParseNumber(am, cur.Fraction)

	if err != nil {
		return ErrInvalidBSONUnmarshal

	}

	*a = fromEquivalentInt(amount, bm.Currency)
	return nil
}

// MarshalBSON is implementation of bson.Marshaller
func (a Money) MarshalBSON() ([]byte, error) {
	currencyCode, amountStr := a.formatAsNumber()

	bm := bsonMoney{
		Amount:   amountStr,
		Currency: currencyCode,
	}

	return bson.Marshal(bm)
}
