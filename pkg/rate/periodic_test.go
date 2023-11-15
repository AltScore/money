package rate

import (
	"github.com/AltScore/money/pkg/money"
	"reflect"
	"testing"
)

func TestPeriodicRate_RoundedByWithPeriod(t *testing.T) {
	type args struct {
		amount money.Money
		period uint
	}
	tests := []struct {
		name string
		rate Periodic
		args args
		want money.Money
	}{
		{
			name: "1% per month",
			rate: NewPeriodicRateFromInt(Month, 42),
			args: args{
				amount: money.FromFloat64(200, "ARS"),
				period: Month,
			},
			want: money.FromFloat64(84, "ARS"),
		},
		{
			name: "120% per year monthly",
			rate: NewPeriodicRateFromFloat64(Year, 120.0),
			args: args{
				amount: money.FromFloat64(4000, "ARS"),
				period: Month,
			},
			want: money.FromFloat64(400, "ARS"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rate.RoundedByWithPeriod(tt.args.amount, tt.args.period); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoundedByWithPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}
