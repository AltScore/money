package parsers

import (
	"strconv"
	"strings"
)

// ParseNumber parses the string for the representation of a number in
// format "dddd.dd" or "dddddd" (without decimal digit)
// Note: If excess decimals contains invalid characters, they are ignored
func ParseNumber(s string, decimals int) (int64, error) {
	dotPos := strings.IndexByte(s, '.')

	var toParse string

	var digits int

	if dotPos < 0 {
		toParse = s
	} else {
		digits = len(s) - dotPos - 1
		if digits > decimals {
			digits = decimals
		}

		firstDigitOutsidePrecision := dotPos + 1 + digits
		toParse = s[:dotPos] + s[dotPos+1:firstDigitOutsidePrecision]

		if firstDigitOutsidePrecision < len(s) && !ArrAllDigits(s[firstDigitOutsidePrecision:]) {
			// If there are invalid characters after the precision, return an error
			return 0, &strconv.NumError{Func: "ParseNumber", Num: s, Err: strconv.ErrSyntax}
		}
	}

	value, err := strconv.ParseInt(toParse, 10, 64)

	if err != nil {
		return 0, err
	}

	for d := digits; d < decimals; d++ {
		value *= 10
	}

	return value, nil
}

func ArrAllDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}
