package parsers

import "strconv"

const zeroes = "0000000000"

func FormatNumber(number int64, decimals int) string {
	s := strconv.FormatInt(number, 10)

	if decimals == 0 {
		return s
	}

	decimalsToAdd := decimals + 1 - len(s)

	if decimalsToAdd > 0 {
		s = zeroes[:decimalsToAdd] + s
	}

	return s[:len(s)-decimals] + "." + s[len(s)-decimals:]
}
