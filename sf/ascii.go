package sf

const (
	lower = ((1 << 26) - 1) << 'a'
	upper = ((1 << 26) - 1) << 'A'
	digit = ((1 << 10) - 1) << '0'
)

func isPrint(c byte) bool { return ' ' <= c && c <= '~' }
func isDigit(c byte) bool { return '0' <= c && c <= '9' }
func isUpper(c byte) bool { return 'A' <= c && c <= 'Z' }
func isLower(c byte) bool { return 'a' <= c && c <= 'z' }
func isAlpha(c byte) bool {
	if isLower(c) {
		return true
	}
	return isUpper(c)
}

func isHexDigit(c uint8) bool {
	if '0' <= c && c <= '9' {
		return true
	}
	if 'A' <= c && c <= 'F' {
		return true
	}
	return 'a' <= c && c <= 'f'
}

func isASCII(c byte, lo, hi uint64) bool {
	var m byte
	if int8(c) >= 0 {
		m = 1
	}
	x := lo
	if c > 63 {
		x = hi
	}
	return uint8(x>>(c%64))&m != 0
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
