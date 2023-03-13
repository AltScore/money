package money

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var centsToNanos int32 = 10000000

func TestFromCommonType(t *testing.T) {
	tests := []struct {
		name string
		args CommonTypeMoney
		want Money
	}{
		{
			name: "zero",
			args: &moneyStub{"MXN", 0, 0},
			want: Zero("MXN"),
		},
		{
			name: "one",
			args: &moneyStub{"MXN", 1, 0},
			want: NewFromInt(1, "MXN"),
		},
		{
			name: "one with nanos",
			args: &moneyStub{"MXN", 1, 1},
			want: NewFromInt(1, "MXN"),
		},
		{
			name: "one with nanos",
			args: &moneyStub{"MXN", 1, 12 * centsToNanos},
			want: MustParse("1.12", "MXN"),
		},
		{
			name: "one with nanos",
			args: &moneyStub{"MXN", -5341, -42 * centsToNanos},
			want: MustParse("-5341.42", "MXN"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FromCommonType(tt.args), "FromCommonType(%v)", tt.args)
		})
	}
}

func TestMoney_AsUnitsAndNanos(t *testing.T) {
	tests := []struct {
		name      string
		m         Money
		wantUnits int64
		wantNanos int32
	}{
		{
			name:      "zero",
			m:         Zero("MXN"),
			wantUnits: 0,
			wantNanos: 0,
		},
		{
			name:      "one",
			m:         NewFromInt(1, "MXN"),
			wantUnits: 1,
			wantNanos: 0,
		},
		{
			name:      "one with nanos",
			m:         MustParse("1.12", "MXN"),
			wantUnits: 1,
			wantNanos: 12 * centsToNanos,
		},
		{
			name:      "one with nanos",
			m:         MustParse("-5341.42", "MXN"),
			wantUnits: -5341,
			wantNanos: -42 * centsToNanos,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.m.AsUnitsAndNanos()
			assert.Equalf(t, tt.wantUnits, got, "AsUnitsAndNanos()")
			assert.Equalf(t, tt.wantNanos, got1, "AsUnitsAndNanos()")
		})
	}
}

func TestMoney_GetUnits(t *testing.T) {
	tests := []struct {
		name string
		m    Money
		want int64
	}{
		{
			name: "zero",
			m:    Zero("MXN"),
			want: 0,
		},
		{
			name: "one",
			m:    NewFromInt(1, "MXN"),
			want: 1,
		},
		{
			name: "one with nanos",
			m:    MustParse("1.12", "MXN"),
			want: 1,
		},
		{
			name: "one with nanos",
			m:    MustParse("-295341.42", "MXN"),
			want: -295341,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.m.GetUnits(), "GetUnits()")
		})
	}
}

func TestMoney_GetNanos(t *testing.T) {
	tests := []struct {
		name string
		m    Money
		want int32
	}{
		{
			name: "zero",
			m:    Zero("MXN"),
			want: 0,
		},
		{
			name: "one",
			m:    NewFromInt(1, "MXN"),
			want: 0,
		},
		{
			name: "one with nanos",
			m:    MustParse("1.12", "MXN"),
			want: 12 * centsToNanos,
		},
		{
			name: "one with nanos",
			m:    MustParse("-5341.42", "MXN"),
			want: -42 * centsToNanos,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.m.GetNanos(), "GetNanos()")
		})
	}
}

type moneyStub struct {
	currencyCode string
	units        int64
	nanos        int32
}

func (m *moneyStub) GetCurrencyCode() string {
	return m.currencyCode
}

func (m *moneyStub) GetUnits() int64 {
	return m.units
}

func (m *moneyStub) GetNanos() int32 {
	return m.nanos
}
