package vm

// decodeSingle decodes a single immediate byte for DUPN/SWAPN.
func decodeSingle(x int) (int, bool) {
	if x > 90 && x < 128 {
		return 0, false
	}
	return (x + 145) % 256, true
}

// decodePair decodes a pair of indices from a single immediate byte for EXCHANGE.
func decodePair(x int) (int, int, bool) {
	if x > 81 && x < 128 {
		return 0, 0, false
	}
	k := x ^ 143
	q := k / 16
	r := k % 16
	if q < r {
		return q + 1, r + 1, true
	}
	return r + 1, 29 - q, true
}
