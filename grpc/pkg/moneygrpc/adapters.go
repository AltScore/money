package moneygrpc

import (
	"github.com/AltScore/money/pkg/money"
	gmoney "google.golang.org/genproto/googleapis/type/money"
)

// MoneyToProto converts money.Money to a google.type.Money proto to be used in gRPC messages
func MoneyToProto(m money.Money) *gmoney.Money {
	units, nanos := m.AsUnitsAndNanos()
	return &gmoney.Money{
		CurrencyCode: m.CurrencyCode(),
		Units:        units,
		Nanos:        nanos,
	}
}
