package money

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidAmount(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "valid amount without decimal point",
			value:   "123",
			wantErr: assert.NoError,
		},
		{
			name:    "valid amount with decimal point",
			value:   "123.42",
			wantErr: assert.NoError,
		},
		{
			name:    "invalid amount with decimal point",
			value:   "123.4.2",
			wantErr: assert.Error,
		},
		{
			name:    "invalid amount with decimal point",
			value:   "123A",
			wantErr: assert.Error,
		},
		{
			name:    "invalid amount with decimal comma",
			value:   "123,42",
			wantErr: assert.Error,
		},
		{
			name:    "invalid amount with spaces prefix",
			value:   " 123.42",
			wantErr: assert.Error,
		},
		{
			name:    "invalid amount with spaces suffix",
			value:   "123.42 ",
			wantErr: assert.Error,
		},
		{
			name:    "invalid amount empty string",
			value:   "",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, IsValidAmount(tt.value), fmt.Sprintf("IsValidAmount(%v)", tt.value))
		})
	}
}
