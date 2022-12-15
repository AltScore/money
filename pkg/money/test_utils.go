package money

import (
	"github.com/stretchr/testify/assert"
)

func EqualAmount(t assert.TestingT, expected Money, actual Money) {
	assert.Equal(t, expected.CurrencyCode(), actual.CurrencyCode())
	comparison, err := expected.TryCmp(actual)
	assert.Nil(t, err, "no errors while comparing %v and %v", actual, expected)
	assert.Equal(t, 0, comparison, "actual %v is not equal to expected %v", actual, expected)
}
