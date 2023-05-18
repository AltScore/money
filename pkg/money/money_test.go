package money

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Currency_can_be_copied(t *testing.T) {
	a, _ := Parse("100", "USD")

	b := a

	assert.Equal(t, a, b)
	assert.NotSame(t, &b, &a)
	assert.Equal(t, a.CurrencyCode(), b.CurrencyCode())
	assert.Equal(t, a.Number(), b.Number())
}

func TestMoney_LessThan(t *testing.T) {
	a := MustParse("100.00", "MXN")
	b := MustParse("200.00", "MXN")
	assert.True(t, a.LessThan(b))
	assert.False(t, b.LessThan(a))
}

func TestMoney_TryCmp(t *testing.T) {
	tests := []struct {
		name    string
		a       Money
		b       Money
		want    int
		wantErr string
	}{
		{name: "a < b", a: MustParse("100.00", "MXN"), b: MustParse("200.00", "MXN"), want: -1},
		{name: "a > b", a: MustParse("200.00", "MXN"), b: MustParse("100.00", "MXN"), want: 1},
		{name: "a == b", a: MustParse("100.00", "MXN"), b: MustParse("100.00", "MXN"), want: 0},
		{name: "different currency", a: MustParse("100.00", "MXN"), b: MustParse("100.00", "ARS"), wantErr: "currencies don't match"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.TryCmp(tt.b)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.Equalf(t, tt.want, got, "TryCmp(%v)", tt.b)
			}
		})
	}
}

func TestNewFromInt(t *testing.T) {
	m := NewFromInt(100, "MXN")
	assert.Equal(t, "MXN", m.CurrencyCode())
	assert.Equal(t, "$100.00", m.String())
}

func TestFromFloat64(t *testing.T) {
	m := FromFloat64(100.5, "MXN")
	assert.Equal(t, "MXN", m.CurrencyCode())
	assert.Equal(t, "$100.50", m.String())
}

func TestMoney_Sub(t *testing.T) {
	tests := []struct {
		name string
		a    Money
		args Money
		want Money
	}{
		{
			name: "Zero - Zero",
			a:    MustParse("0.00", "MXN"),
			args: MustParse("0.00", "MXN"),
			want: MustParse("0.00", "MXN"),
		},
		{
			name: "Zero - Positive",
			a:    MustParse("0.00", "MXN"),
			args: MustParse("100.00", "MXN"),
			want: MustParse("-100.00", "MXN"),
		},
		{
			name: "Zero - Negative",
			a:    MustParse("0.00", "MXN"),
			args: MustParse("-100.00", "MXN"),
			want: MustParse("100.00", "MXN"),
		},
		{
			name: "Positive - Zero",
			a:    MustParse("100.00", "MXN"),
			args: MustParse("0.00", "MXN"),
			want: MustParse("100.00", "MXN"),
		},
		{
			name: "Positive - Positive",

			a:    MustParse("100.00", "MXN"),
			args: MustParse("100.00", "MXN"),
			want: MustParse("0.00", "MXN"),
		},
		{
			name: "Positive - Negative",
			a:    MustParse("100.00", "MXN"),
			args: MustParse("-100.00", "MXN"),
			want: MustParse("200.00", "MXN"),
		},
		{
			name: "Big Positive - Small Positive",

			a:    MustParse("1000.00", "MXN"),
			args: MustParse("200.00", "MXN"),
			want: MustParse("800.00", "MXN"),
		},
		{
			name: "Small Positive - Big Positive",

			a:    MustParse("100.00", "MXN"),
			args: MustParse("2000.00", "MXN"),
			want: MustParse("-1900.00", "MXN"),
		},
		{
			name: "Negative - Zero",
			a:    MustParse("-100.00", "MXN"),
			args: MustParse("0.00", "MXN"),
			want: MustParse("-100.00", "MXN"),
		},
		{
			name: "Negative - Zero different currency",
			a:    MustParse("-100.00", "MXN"),
			args: MustParse("0.00", "USD"),
			want: MustParse("-100.00", "MXN"),
		},
		{
			name: "Positive - Zero different currency",
			a:    MustParse("0", "MXN"),
			args: MustParse("-79.12", "USD"),
			want: MustParse("79.12", "USD"),
		},
		{
			name: "Zero - Zero different currency",
			a:    MustParse("0", "MXN"),
			args: MustParse("0", "USD"),
			want: MustParse("0", "MXN"),
		},
		{
			name: "Zero - Positive different currency",
			a:    MustParse("0", "MXN"),
			args: MustParse("65.87", "USD"),
			want: MustParse("-65.87", "USD"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.a.Sub(tt.args), "Sub(%v)", tt.args)
		})
	}
}

func TestMoney_Negated(t *testing.T) {
	tests := []struct {
		name string
		a    Money
		want Money
	}{
		{
			name: "Zero",
			a:    MustParse("0.00", "MXN"),
			want: MustParse("0.00", "MXN"),
		},
		{
			name: "Positive",
			a:    MustParse("100.00", "MXN"),
			want: MustParse("-100.00", "MXN"),
		},
		{
			name: "Negative",
			a:    MustParse("-100.00", "MXN"),
			want: MustParse("100.00", "MXN"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.a.Negated(), "Negated()")
		})
	}
}

func TestMoney_IsGreaterThan(t *testing.T) {
	tests := []struct {
		name string
		a    Money
		args Money
		want bool
	}{
		{
			name: "Zero > Zero",
			a:    MustParse("0.00", "MXN"),
			args: MustParse("0.00", "MXN"),
			want: false,
		},
		{
			name: "Zero > Positive",
			a:    MustParse("0.00", "MXN"),
			args: MustParse("100.00", "MXN"),
			want: false,
		},
		{
			name: "Zero > Negative",
			a:    MustParse("0.00", "MXN"),
			args: MustParse("-100.00", "MXN"),
			want: true,
		},
		{
			name: "Positive > Zero",
			a:    MustParse("100.00", "MXN"),
			args: MustParse("0.00", "MXN"),
			want: true,
		},
		{
			name: "Positive > Positive",
			a:    MustParse("100.00", "MXN"),
			args: MustParse("100.00", "MXN"),
			want: false,
		},
		{
			name: "Positive > Negative",
			a:    MustParse("100.00", "MXN"),
			args: MustParse("-100.00", "MXN"),
			want: true,
		},
		{
			name: "Negative > Zero",
			a:    MustParse("-100.00", "MXN"),
			args: MustParse("0.00", "MXN"),
			want: false,
		},
		{
			name: "Negative > Positive",
			a:    MustParse("-100.00", "MXN"),
			args: MustParse("100.00", "MXN"),
			want: false,
		},
		{
			name: "Negative > Negative",
			a:    MustParse("-100.00", "MXN"),
			args: MustParse("-100.00", "MXN"),
			want: false,
		},
		{
			name: "Zero > Positive different currency",
			a:    MustParse("0.00", "USD"),
			args: MustParse("100.00", "MXN"),
			want: false,
		},
		{
			name: "Positive > Zero different currency",
			a:    MustParse("125", "USD"),
			args: MustParse("0.00", "MXN"),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.a.IsGreaterThan(tt.args), "IsGreaterThan(%v)", tt.args)
		})
	}
}

func TestMoney_RoundDiv(t *testing.T) {
	tests := []struct {
		name    string
		a       Money
		divider int64
		want    Money
	}{
		{
			name:    "Zero",
			a:       MustParse("0.00", "MXN"),
			divider: 1,
			want:    MustParse("0.00", "MXN"),
		},
		{
			name:    "Positive",
			a:       MustParse("100.00", "MXN"),
			divider: 2,
			want:    MustParse("50.00", "MXN"),
		},
		{
			name:    "Rounded 0.01",
			a:       MustParse("0.01", "MXN"),
			divider: 2,
			want:    MustParse("0.00", "MXN"), // because Half Even Rounding
		},
		{
			name:    "Rounded 0.05",
			a:       MustParse("0.05", "MXN"),
			divider: 2,
			want:    MustParse("0.02", "MXN"), // because Half Even Rounding
		},
		{
			name:    "Rounded 0.50 by 20",
			a:       MustParse("0.50", "MXN"),
			divider: 20,
			want:    MustParse("0.02", "MXN"), // because Half Even Rounding
		},
		{
			name:    "Rounded 0.51 by 20",
			a:       MustParse("0.51", "MXN"),
			divider: 20,
			want:    MustParse("0.03", "MXN"), // because Half Even Rounding
		},
		{
			name:    "Rounded 0.49 by 20",
			a:       MustParse("0.49", "MXN"),
			divider: 20,
			want:    MustParse("0.02", "MXN"),
		},
		{
			name:    "Rounded 0.22 by 4",
			a:       MustParse("0.22", "MXN"),
			divider: 6,
			want:    MustParse("0.04", "MXN"), // because Half Even Rounding
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.a.RoundedDiv(tt.divider), "RoundedDiv(%v)", tt.divider)
		})
	}
}

func TestMoney_RoundDiv_half_even_rounding_mode(t *testing.T) {
	tests := []struct {
		name    string
		a       Money
		divider int64
		want    Money
	}{
		{
			name:    "Zero",
			a:       MustParse("0.05", "MXN"),
			divider: 10,
			want:    MustParse("0.00", "MXN"),
		},
		{
			name:    "0.15",
			a:       MustParse("0.15", "MXN"),
			divider: 10,
			want:    MustParse("0.02", "MXN"),
		},
		{
			name:    "0.25",
			a:       MustParse("0.25", "MXN"),
			divider: 10,
			want:    MustParse("0.02", "MXN"),
		},
		{
			name:    "0.35",
			a:       MustParse("0.35", "MXN"),
			divider: 10,
			want:    MustParse("0.04", "MXN"),
		},
		{
			name:    "0.45",
			a:       MustParse("0.45", "MXN"),
			divider: 10,
			want:    MustParse("0.04", "MXN"),
		},
		{
			name:    "0.55",
			a:       MustParse("0.55", "MXN"),
			divider: 10,
			want:    MustParse("0.06", "MXN"),
		},
		{
			name:    "0.65",
			a:       MustParse("0.65", "MXN"),
			divider: 10,
			want:    MustParse("0.06", "MXN"),
		},
		{
			name:    ".075",
			a:       MustParse("0.75", "MXN"),
			divider: 10,
			want:    MustParse("0.08", "MXN"),
		},
		{
			name:    "0.85",
			a:       MustParse("0.85", "MXN"),
			divider: 10,
			want:    MustParse("0.08", "MXN"),
		},
		{
			name:    "0.95",
			a:       MustParse("0.95", "MXN"),
			divider: 10,
			want:    MustParse("0.10", "MXN"),
		},

		{
			name:    "-0.15",
			a:       MustParse("-0.15", "MXN"),
			divider: 10,
			want:    MustParse("-0.02", "MXN"),
		},
		{
			name:    "-0.25",
			a:       MustParse("-0.25", "MXN"),
			divider: 10,
			want:    MustParse("-0.02", "MXN"),
		},
		{
			name:    "-0.35",
			a:       MustParse("-0.35", "MXN"),
			divider: 10,
			want:    MustParse("-0.04", "MXN"),
		},
		{
			name:    "-0.45",
			a:       MustParse("-0.45", "MXN"),
			divider: 10,
			want:    MustParse("-0.04", "MXN"),
		},
		{
			name:    "-0.55",
			a:       MustParse("-0.55", "MXN"),
			divider: 10,
			want:    MustParse("-0.06", "MXN"),
		},
		{
			name:    "-0.65",
			a:       MustParse("-0.65", "MXN"),
			divider: 10,
			want:    MustParse("-0.06", "MXN"),
		},
		{
			name:    "-.075",
			a:       MustParse("-0.75", "MXN"),
			divider: 10,
			want:    MustParse("-0.08", "MXN"),
		},
		{
			name:    "-0.85",
			a:       MustParse("-0.85", "MXN"),
			divider: 10,
			want:    MustParse("-0.08", "MXN"),
		},
		{
			name:    "-0.95",
			a:       MustParse("-0.95", "MXN"),
			divider: 10,
			want:    MustParse("-0.10", "MXN"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.a.RoundedDiv(tt.divider), "%v.RoundedDiv(%v)", tt.a, tt.divider)
		})
	}
}

func TestMoney_String(t *testing.T) {
	tests := []struct {
		name string
		a    Money
		want string
	}{
		{
			name: "Zero",
			a:    MustParse("0.00", "MXN"),
			want: "$0.00",
		},
		{
			name: "Negative",
			a:    MustParse("-1.00", "MXN"),
			want: "-$1.00",
		},
		{
			name: "Positive",
			a:    MustParse("1.00", "MXN"),
			want: "$1.00",
		},
		{
			name: "Negative with cents",
			a:    MustParse("-1.01", "MXN"),
			want: "-$1.01",
		},
		{
			name: "Positive with cents",
			a:    MustParse("1.01", "MXN"),
			want: "$1.01",
		},
		{
			name: "empty",
			a:    Money{},
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.a.String(), "String()")
		})
	}
}

func TestMoney_GoString(t *testing.T) {
	tests := []struct {
		name string
		a    Money
		want string
	}{
		{
			name: "Zero",
			a:    MustParse("0.00", "MXN"),
			want: `money.FromFloat64(0, "MXN")`,
		},
		{
			name: "Negative",
			a:    MustParse("-1.00", "MXN"),
			want: `money.FromFloat64(-1, "MXN")`,
		},
		{
			name: "Positive",
			a:    MustParse("1.00", "MXN"),
			want: `money.FromFloat64(1, "MXN")`,
		},
		{
			name: "Negative with cents",
			a:    MustParse("-1.01", "MXN"),
			want: `money.FromFloat64(-1.01, "MXN")`,
		},
		{
			name: "Positive with cents",
			a:    MustParse("1.01", "MXN"),
			want: `money.FromFloat64(1.01, "MXN")`,
		},
		{
			name: "empty",
			a:    Money{},
			want: `money.FromFloat64(0, "")`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.a.GoString(), "GoString()")
		})
	}
}

func TestMoney_Sprintf(t *testing.T) {
	tests := []struct {
		name string
		a    Money
		want string
	}{
		{
			name: "Zero",
			a:    MustParse("0.00", "MXN"),
			want: "$0.00",
		},
		{
			name: "Negative",
			a:    MustParse("-1.00", "MXN"),
			want: "-$1.00",
		},
		{
			name: "Positive",
			a:    MustParse("1.00", "MXN"),
			want: "$1.00",
		},
		{
			name: "Negative with cents",
			a:    MustParse("-1.01", "MXN"),
			want: "-$1.01",
		},
		{
			name: "Positive with cents",
			a:    MustParse("1.01", "MXN"),
			want: "$1.01",
		},
		{
			name: "empty",
			a:    Money{},
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, fmt.Sprintf("%s", tt.a), "Sprintf(%%s)")
			assert.Equalf(t, tt.want, fmt.Sprintf("%v", tt.a), "Sprintf(%%v)")
		})
	}
}

func TestMoney_TryAdd(t *testing.T) {

	tests := []struct {
		name    string
		value   Money
		args    Money
		want    Money
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Add 1.00 + 1.00",
			value:   MustParse("1.00", "MXN"),
			args:    MustParse("1.00", "MXN"),
			want:    MustParse("2.00", "MXN"),
			wantErr: assert.NoError,
		},
		{
			name:    "Add 0.00 + 42.23",
			value:   MustParse("0.00", "MXN"),
			args:    MustParse("42.23", "MXN"),
			want:    MustParse("42.23", "MXN"),
			wantErr: assert.NoError,
		},
		{
			name:    "Add 1.00 + empty",
			value:   MustParse("1.00", "MXN"),
			args:    Money{},
			want:    MustParse("1.00", "MXN"),
			wantErr: assert.NoError,
		},
		{
			name:    "Add empty + 1.00",
			value:   Money{},
			args:    MustParse("1.00", "MXN"),
			want:    MustParse("1.00", "MXN"),
			wantErr: assert.NoError,
		},
		{
			name:    "Add 0 + 3132.23",
			value:   MustParse("0", "MXN"),
			args:    MustParse("3132.23", "MXN"),
			want:    MustParse("3132.23", "MXN"),
			wantErr: assert.NoError,
		},
		{
			name:    "Add 54343.12 + 0",
			value:   MustParse("54343.12", "MXN"),
			args:    MustParse("0", "MXN"),
			want:    MustParse("54343.12", "MXN"),
			wantErr: assert.NoError,
		},
		{
			name:    "Add 0 + 0",
			value:   MustParse("0", "MXN"),
			args:    MustParse("0", "MXN"),
			want:    MustParse("0", "MXN"),
			wantErr: assert.NoError,
		},
		{
			name:    "Add 0 + empty",
			value:   MustParse("0", "MXN"),
			args:    Money{},
			want:    MustParse("0", "MXN"),
			wantErr: assert.NoError,
		},
		{
			name:    "Add empty + 0",
			value:   Money{},
			args:    MustParse("0", "MXN"),
			want:    MustParse("0", "MXN"),
			wantErr: assert.NoError,
		},
		{
			name:    "Add empty + empty",
			value:   Money{},
			args:    Money{},
			want:    Money{},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.value
			got, err := a.TryAdd(tt.args)
			if !tt.wantErr(t, err, fmt.Sprintf("TryAdd(%v)", tt.args)) {
				return
			}
			assert.Equalf(t, tt.want, got, "TryAdd(%v)", tt.args)
		})
	}
}
