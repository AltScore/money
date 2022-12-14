package utils

func HalfEvenRounding(a, b int64) int64 {
	if a > 0 && b < 0 || a < 0 && b > 0 {
		return -HalfEvenRounding(-a, b)
	}

	r := a * 2 / b // make division keeping first bit of remaining

	switch r & 0b11 {
	case 0b11:
		// odd number, with excess
		return r/2 + 1

	case 0b10:
		// odd number, without excess
		return r / 2

	case 0b01:
		// even number, with excess
		if r*b == a*2 {
			return r / 2
		}

		return r/2 + 1

	default:
		// 0b00
		// even number, without excess
		return r / 2
	}
}
