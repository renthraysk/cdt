package sf

func isLower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// isASCIIValid checks whether all characters in the string s are valid ASCII characters
// which are in the 128 bitset lo and hi.
func isASCIIValid(s string, lo, hi uint64) bool {
	const maxASCII = '\x7F'

	var h, l uint64
	i := 0
	for ; i < len(s) && s[i] <= maxASCII; i++ {
		x, y := uint64(s[i]/64), s[i]%64 // x ∈ {0, 1}
		h |= x << y
		l |= (x ^ 1) << y
	}
	return i >= len(s) && (l&^lo)|(h&^hi) == 0
}
