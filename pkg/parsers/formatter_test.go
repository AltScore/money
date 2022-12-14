package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatNumber(t *testing.T) {
	type args struct {
		number   int64
		decimals int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"zero", args{number: 0, decimals: 0}, "0"},
		{"zero with decimals", args{number: 0, decimals: 2}, "0.00"},
		{"number", args{number: 42, decimals: 0}, "42"},
		{"one digit", args{number: 2, decimals: 2}, "0.02"},
		{"two digit", args{number: 42, decimals: 2}, "0.42"},
		{"three digit", args{number: 154, decimals: 2}, "1.54"},
		{"big number with decimals", args{number: 78987546, decimals: 3}, "78987.546"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FormatNumber(tt.args.number, tt.args.decimals), "FormatNumber(%v, %v)", tt.args.number, tt.args.decimals)
		})
	}
}
