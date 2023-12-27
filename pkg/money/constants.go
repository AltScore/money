package money

import (
	"fmt"
	"regexp"
)

const (
	AmountAsStringFormat = `^(?:-)?\d+(\.\d+)?$`
)

var AmountAsStringRegexp = regexp.MustCompile(AmountAsStringFormat)

// IsValidAmount returns no error if the amount is valid (positive or negative), an error if it is not.
func IsValidAmount(value string) error {
	if AmountAsStringRegexp.MatchString(value) {
		return nil
	}
	return fmt.Errorf("invalid amount: %s", value)
}
