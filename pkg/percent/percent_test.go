package percent

import (
	"fmt"
	"testing"

	"github.com/AltScore/money/pkg/money"
	"github.com/stretchr/testify/assert"
)

func Test_Percent_can_be_compared_with_operator(t *testing.T) {
	// nolint:staticcheck // this comparison is just the test
	assert.True(t, Percent(42) == Percent(42))
	assert.False(t, Percent(42) == Percent(24))
	assert.True(t, Percent(42) != Percent(24))
	assert.True(t, Percent(42) > Percent(24))
	assert.False(t, Percent(42) < Percent(24))

}

func TestFromFloat64(t *testing.T) {

	tests := []struct {
		name string
		pct  float64
		want Percent
	}{
		{"Zero", 0, MustParse("0.00")},
		{"One", 100, MustParse("100")},
		{"few decimals", 42.25, MustParse("42.25")},
		{"many decimals rounding down", 67.45635, MustParse("67.4564")},
		{"many decimals rounding up", 125.416725, MustParse("125.4167")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FromFloat64(tt.pct), "FromFloat64(%v)", tt.pct)
		})
	}
}

func TestFromFactor(t *testing.T) {

	tests := []struct {
		name string
		pct  float64
		want Percent
	}{
		{"Zero", 0, MustParse("0.00")},
		{"One", 1, MustParse("100")},
		{"few decimals", 0.4225, MustParse("42.25")},
		{"many decimals rounding down", 0.6745635, MustParse("67.4564")},
		{"many decimals rounding up", 1.25416725, MustParse("125.4167")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FromFactor(tt.pct), "FromFloat64(%v)", tt.pct)
		})
	}
}

func TestFromFraction64(t *testing.T) {
	type args struct {
		partial float64
		total   float64
	}
	tests := []struct {
		name string
		args args
		want Percent
	}{
		{"Zero", args{partial: 0, total: 0}, MustParse("0.00")},
		{"partial lower", args{partial: 11, total: 33}, MustParse("33.3333")},
		{"partial higher", args{partial: 125, total: 25}, MustParse("500.00")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FromFraction64(tt.args.partial, tt.args.total), "FromFraction64(%v, %v)", tt.args.partial, tt.args.total)
		})
	}
}

func TestFromFraction(t *testing.T) {
	type args struct {
		partial money.Money
		total   money.Money
	}
	tests := []struct {
		name    string
		args    args
		want    Percent
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Zero",
			args: args{
				partial: money.Zero("ARS"),
				total:   money.Zero("ARS"),
			},
			want:    MustParse("0.00"),
			wantErr: isDivisionByZeroError,
		},
		{
			name: "Different currencies",
			args: args{
				partial: money.MustParse("100", "ARS"),
				total:   money.MustParse("20", "MXN"),
			},
			want:    MustParse("0.00"),
			wantErr: isDifferentCurrenciesError,
		},
		{
			name: "Good percent",
			args: args{
				partial: money.MustParse("20", "ARS"),
				total:   money.MustParse("100", "ARS"),
			},
			want:    MustParse("20.00"),
			wantErr: assert.NoError,
		},
		{
			name: "Zero partial",
			args: args{
				partial: money.Zero("ARS"),
				total:   money.MustParse("100", "ARS"),
			},
			want:    MustParse("0.00"),
			wantErr: assert.NoError,
		},
		{
			name: "Zero total",
			args: args{
				partial: money.MustParse("100", "ARS"),
				total:   money.Zero("ARS"),
			},
			want:    MustParse("0.00"),
			wantErr: isDivisionByZeroError,
		},
		{
			name: "Zero partial different currency",
			args: args{
				partial: money.Zero("MXN"),
				total:   money.MustParse("100", "ARS"),
			},
			want:    MustParse("0.00"),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromFraction(tt.args.partial, tt.args.total)
			if !tt.wantErr(t, err, fmt.Sprintf("FromFraction(%v, %v)", tt.args.partial, tt.args.total)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FromFraction(%v, %v)", tt.args.partial, tt.args.total)
		})
	}
}

func TestChangePeriod(t *testing.T) {
	type args struct {
		rate          string
		periodSize    int
		newPeriodSize int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "same period",
			args: args{rate: "3.5", periodSize: 29, newPeriodSize: 29},
			want: "3.5",
		}, {
			name: "monthly rate to bimonthly",
			args: args{rate: "3.5", periodSize: 30, newPeriodSize: 60},
			want: "7.1225",
		},
		{
			name: "monthly rate to quarter monthly",
			args: args{rate: "3.5", periodSize: 30, newPeriodSize: 90},
			want: "10.8718",
		},
		{
			name: "monthly rate to biweekly",
			args: args{rate: "3.5", periodSize: 30, newPeriodSize: 15},
			want: "1.7349",
		},
		{
			name: "monthly rate to 3 weeks",
			args: args{rate: "3.5", periodSize: 30, newPeriodSize: 45},
			want: "5.2957",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate := MustParse(tt.args.rate)
			got := rate.ChangePeriod(tt.args.periodSize, tt.args.newPeriodSize)
			assert.Equal(t, MustParse(tt.want), got)
		})
	}
}

func TestPercent_By(t *testing.T) {
	tests := []struct {
		name    string
		percent Percent
		amount  money.Money
		want    money.Money
	}{
		{
			name:    "Zero",
			percent: MustParse("0.00"),
			amount:  money.Zero("ARS"),
			want:    money.Zero("ARS"),
		},
		{
			name:    "Zero percent",
			percent: MustParse("0.00"),
			amount:  money.MustParse("20", "MXN"),
			want:    money.Zero("MXN"),
		},
		{
			name:    "Zero amount",
			percent: MustParse("20"),
			amount:  money.MustParse("0", "ARS"),
			want:    money.MustParse("0", "ARS"),
		},
		{
			name:    "10%",
			percent: MustParse("10"),
			amount:  money.MustParse("1234", "ARS"),
			want:    money.MustParse("123.4", "ARS"),
		},
		{
			name:    "20%",
			percent: MustParse("20"),
			amount:  money.MustParse("61.20", "MXN"),
			want:    money.MustParse("12.24", "MXN"),
		},
		{
			name:    "16%",
			percent: MustParse("16"),
			amount:  money.MustParse("1495.41", "MXN"),
			want:    money.MustParse("239.26", "MXN"), // Exact value is 239.26560
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.percent.By(tt.amount)

			assert.Equalf(t, tt.want, got, "%v.By(%v)", tt.percent, tt.amount)
		})
	}
}

func TestPercent_RoundedBy(t *testing.T) {
	tests := []struct {
		name    string
		percent Percent
		amount  money.Money
		want    money.Money
	}{
		{
			name:    "Zero",
			percent: MustParse("0.00"),
			amount:  money.Zero("ARS"),
			want:    money.Zero("ARS"),
		},
		{
			name:    "Zero percent",
			percent: MustParse("0.00"),
			amount:  money.MustParse("20", "MXN"),
			want:    money.Zero("MXN"),
		},
		{
			name:    "Zero amount",
			percent: MustParse("20"),
			amount:  money.MustParse("0", "ARS"),
			want:    money.MustParse("0", "ARS"),
		},
		{
			name:    "10%",
			percent: MustParse("10"),
			amount:  money.MustParse("1234", "ARS"),
			want:    money.MustParse("123.4", "ARS"),
		},
		{
			name:    "20%",
			percent: MustParse("20"),
			amount:  money.MustParse("61.20", "MXN"),
			want:    money.MustParse("12.24", "MXN"),
		},
		{
			name:    "16%",
			percent: MustParse("16"),
			amount:  money.MustParse("1495.41", "MXN"),
			want:    money.MustParse("239.27", "MXN"), // Exact value is 239.26560
		},
		{
			name:    "16% of 19.48",
			percent: MustParse("16"),
			amount:  money.MustParse("19.48", "MXN"),
			want:    money.MustParse("3.12", "MXN"), // Exact value is 239.26560
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.percent.RoundedBy(tt.amount)

			assert.Equalf(t, tt.want, got, "%v.By(%v)", tt.percent, tt.amount)
		})
	}
}

func isDivisionByZeroError(t assert.TestingT, err error, args ...interface{}) bool {
	return assert.EqualError(t, err, "division by zero", args...)
}

func isDifferentCurrenciesError(t assert.TestingT, err error, args ...interface{}) bool {
	return assert.EqualError(t, err, "currencies don't match", args...)
}

func TestPercent_String(t *testing.T) {
	tests := []struct {
		name string
		p    Percent
		want string
	}{
		{
			name: "Zero",
			p:    MustParse("0.00"),
			want: "0",
		},
		{
			name: "10%",
			p:    MustParse("10.00"),
			want: "10",
		},
		{
			name: "10.5%",
			p:    MustParse("10.50"),
			want: "10.5",
		},
		{
			name: "10.05%",
			p:    MustParse("10.05"),
			want: "10.05",
		},
		{
			name: "10.005%",
			p:    MustParse("10.005"),
			want: "10.005",
		},
		{
			name: "42%",
			p:    FromFactor(0.42),
			want: "42",
		},
		{
			name: "7.42",
			p:    FromFloat64(7.42),
			want: "7.42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.p.String(), "String()")
		})
	}
}

func TestPercent_ExtractPercentFromTotal(t *testing.T) {
	tests := []struct {
		name  string
		p     Percent
		total money.Money
		want  money.Money
	}{
		{
			name:  "Zero percent",
			p:     MustParse("0.00"),
			total: money.MustParse("100", "USD"),
			want:  money.MustParse("0", "USD"),
		},
		{
			name:  "Zero amount",
			p:     MustParse("10.00"),
			total: money.MustParse("0", "USD"),
			want:  money.MustParse("0", "USD"),
		},
		{
			name:  "10%",
			p:     MustParse("23.00"),
			total: money.MustParse("123", "USD"),
			want:  money.MustParse("23", "USD"),
		},
		{
			name:  "10.5%",
			p:     MustParse("10.50"),
			total: money.MustParse("61.20", "MXN"),
			want:  money.MustParse("5.81", "MXN"),
		},
		{
			name:  "3% big amount",
			p:     MustParse("3"),
			total: money.MustParse("13500", "MXN"),
			want:  money.MustParse("393.20", "MXN"),
		},
		{
			name:  "Inexact percent",
			p:     MustParse("16"),
			total: money.MustParse("1734.68", "MXN"),
			want:  money.MustParse("239.26", "MXN"), // Exact is 239.26560
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.p.ExtractPercentFromTotal(tt.total), "ExtractPercentFromTotal(%v)", tt.total)
		})
	}
}

func TestPercent_ExtractRoundedPercentFromTotal(t *testing.T) {
	tests := []struct {
		name  string
		p     Percent
		total money.Money
		want  money.Money
	}{
		{
			name:  "Zero percent",
			p:     MustParse("0.00"),
			total: money.MustParse("100", "USD"),
			want:  money.MustParse("0", "USD"),
		},
		{
			name:  "Zero amount",
			p:     MustParse("10.00"),
			total: money.MustParse("0", "USD"),
			want:  money.MustParse("0", "USD"),
		},
		{
			name:  "10%",
			p:     MustParse("23.00"),
			total: money.MustParse("123", "USD"),
			want:  money.MustParse("23", "USD"),
		},
		{
			name:  "10.5%",
			p:     MustParse("10.50"),
			total: money.MustParse("61.20", "MXN"),
			want:  money.MustParse("5.82", "MXN"), // 5.815384615
		},
		{
			name:  "3% big amount",
			p:     MustParse("3"),
			total: money.MustParse("13500", "MXN"),
			want:  money.MustParse("393.20", "MXN"),
		},
		{
			name:  "Inexact percent",
			p:     MustParse("16"),
			total: money.MustParse("1734.68", "MXN"),
			want:  money.MustParse("239.27", "MXN"), // Exact is 239.26560
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.p.ExtractRoundedPercentFromTotal(tt.total), "ExtractPercentFromTotal(%v)", tt.total)
		})
	}
}
