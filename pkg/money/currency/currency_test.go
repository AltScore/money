package currency

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		name         string
		currencyCode string
		want         bool
	}{
		{
			name:         "valid currency code ARS",
			currencyCode: "ARS",
			want:         true,
		},
		{
			name:         "valid currency code USD",
			currencyCode: "USD",
			want:         true,
		},
		{
			name:         "valid currency code MXN",
			currencyCode: "MXN",
			want:         true,
		},
		{
			name:         "invalid currency code",
			currencyCode: "ABC",
			want:         false,
		},
		{
			name:         "invalid empty currency code",
			currencyCode: "",
			want:         false,
		},
		{
			name:         "invalid short currency code",
			currencyCode: "AR",
			want:         false,
		},
		{
			name:         "invalid long currency code",
			currencyCode: "ARS2",
			want:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.currencyCode); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name         string
		currencyCode string
		wantErr      string
	}{
		{
			name:         "valid currency code ARS",
			currencyCode: "ARS",
			wantErr:      "",
		},
		{
			name:         "valid currency code USD",
			currencyCode: "USD",
			wantErr:      "",
		},
		{
			name:         "valid currency code MXN",
			currencyCode: "MXN",
			wantErr:      "",
		},
		{
			name:         "invalid currency code",
			currencyCode: "ABC",
			wantErr:      "invalid currency code: ABC",
		},
		{
			name:         "invalid empty currency code",
			currencyCode: "",
			wantErr:      "invalid currency code: ",
		},
		{
			name:         "invalid short currency code",
			currencyCode: "AR",
			wantErr:      "invalid currency code: AR",
		},
		{
			name:         "invalid long currency code",
			currencyCode: "ARS2",
			wantErr:      "invalid currency code: ARS2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Check(tt.currencyCode)

			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.wantErr)

			}
		})
	}
}
