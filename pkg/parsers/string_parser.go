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

		toParse = s[:dotPos] + s[dotPos+1:dotPos+1+digits]
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
