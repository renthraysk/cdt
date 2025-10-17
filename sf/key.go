package sf

import "slices"

func isLower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// https://www.rfc-editor.org/rfc/rfc8941.html#name-serializing-a-key
func isKeyValid(key string) bool {
	const (
		lcalpha = ((1 << 26) - 1) << 'a'
		digit   = ((1 << 10) - 1) << '0'
		set     = lcalpha | digit | 1<<'_' | 1<<'-' | 1<<'.' | 1<<'*'
	)
	return len(key) > 0 &&
		(isLower(key[0]) || key[0] == '*') &&
		isASCIIValid(key[1:], set%(1<<64), set>>64)
}

func AppendKeyString(p []byte, key, value string) ([]byte, bool) {
	if !isKeyValid(key) {
		return p, false
	}
	n, ok := stringCountEscapeChars(value)
	if !ok {
		return p, false
	}
	q := slices.Grow(p, len(key)+len(`="`)+len(value)+n+len(`"`))
	q = append(q, key...)
	q = append(q, '=', '"')
	if n == 0 {
		q = append(q, value...)
		return append(q, '"'), true
	}
	if q, ok = appendStringEscape(q, value); !ok {
		return p, false
	}
	return append(q, '"'), false
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
