package moneygrpc

import (
	"github.com/AltScore/money/pkg/money"
	gmoney "google.golang.org/genproto/googleapis/type/money"

	"reflect"
	"testing"
)

func TestMoneyToProto(t *testing.T) {
	tests := []struct {
		name string
		args money.Money
		want *gmoney.Money
	}{
		{
			name: "converts money to proto",
			args: money.NewFromInt(100, "USD"),
			want: &gmoney.Money{
				CurrencyCode: "USD",
				Units:        100,
				Nanos:        0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MoneyToProto(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MoneyToProto() = %v, want %v", got, tt.want)
			}
		})
	}
}
