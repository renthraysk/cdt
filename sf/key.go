package sf

import "slices"

// https://www.rfc-editor.org/rfc/rfc8941.html#name-serializing-a-key
func KeyValid(key string) bool {
	const keySet = lower | digit | 1<<'_' | 1<<'-' | 1<<'.' | 1<<'*'

	return len(key) > 0 &&
		(isLower(key[0]) || key[0] == '*') &&
		isASCIIValid(key[1:], keySet%(1<<64), keySet>>64)
}

func KeyAppendString(p []byte, key, value string) ([]byte, bool) {
	if !KeyValid(key) {
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
	} else if q, ok = stringAppendEscape(q, value); !ok {
		return p, false
	}
	return append(q, '"'), true
}
