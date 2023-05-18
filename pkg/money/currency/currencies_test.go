package currency

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddCurrency(t *testing.T) {
	// WHEN we add a currency
	got := AddCurrency("XYZ", "#", "$1", ":", "/", 4)

	// THEN we expect the currency to be added
	find := Get("XYZ")

	require.NotNil(t, find)
	require.Equal(t, got, find)

	require.Equal(t, "XYZ", find.Code)
	require.Equal(t, "#", find.Grapheme)
	require.Equal(t, "$1", find.Template)
	require.Equal(t, ":", find.Decimal)
	require.Equal(t, "/", find.Thousand)
	require.Equal(t, 4, find.Fraction)
}
