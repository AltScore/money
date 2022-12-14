package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseAmount(t *testing.T) {
	type args struct {
		s        string
		decimals int
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr string
	}{
		{name: "zero", args: args{"0", 2}, want: 0, wantErr: ""},
		{name: "integer, no digits", args: args{"1", 0}, want: 1, wantErr: ""},
		{name: "integer, digits", args: args{"1", 2}, want: 100, wantErr: ""},
		{name: "float, no digits", args: args{"1.23", 0}, want: 1, wantErr: ""},
		{name: "float, few digits", args: args{"1.2", 3}, want: 1200, wantErr: ""},
		{name: "float, exact digits", args: args{"1.23", 2}, want: 123, wantErr: ""},
		{name: "float, too any digits", args: args{"1.2345", 2}, want: 123, wantErr: ""},
		{name: "negative float, too any digits", args: args{"-1.2345", 2}, want: -123, wantErr: ""},
		{name: "erroneous value", args: args{"-123X45", 2}, want: 0, wantErr: "invalid syntax"},
		{name: "ignore characters in excess decimals", args: args{"-1.23X45", 2}, want: -123, wantErr: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNumber(tt.args.s, tt.args.decimals)

			if tt.wantErr == "" {
				if !assert.Nil(t, err) {
					return
				}
			} else {
				if !assert.ErrorContains(t, err, tt.wantErr) {
					return
				}
			}
			assert.Equalf(t, tt.want, got, "ParseNumber(%v, %v)", tt.args.s, tt.args.decimals)
		})
	}
}
