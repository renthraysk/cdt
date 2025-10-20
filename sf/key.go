package sf

import "slices"

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
	return append(q, '"'), true
}
