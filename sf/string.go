package sf

import (
	"slices"
)

// https://www.rfc-editor.org/rfc/rfc9651.html#name-strings

func String(v []string) (string, bool) {
	if len(v) != 1 {
		return "", false
	}
	s, r, ok := stringCut(v[0])
	if !ok || len(s) < 2 || len(r) != 0 {
		return "", false
	}
	return stringUnescape(s[1 : len(s)-1])
}

func stringUnescape(s string) (string, bool) {
	// check if valid and if any unescaping to do
	i := 0
	for i < len(s) && isPrint(s[i]) && s[i] != '\\' {
		i++
	}
	if i >= len(s) {
		return s, true
	}
	if !isPrint(s[i]) {
		return "", false
	}
	// s[i] == '\\'
	if j := i + len(`\`); j >= len(s) || (s[j] != '\\' && s[j] != '"') {
		return "", false
	}
	// have a valid escape sequence

	// try and avoid an allocation for short <64b strings
	dst := slices.Grow(make([]byte, 0, 64), len(s))
	dst = append(dst, s[:i]...)
	dst, ok := stringAppendUnescape(dst, s[i:])
	if !ok {
		return "", false
	}
	return string(dst), true
}

func stringCut(s string) (string, string, bool) {
	if len(s) <= 0 || s[0] != '"' {
		return "", s, false
	}
	i := 1
	for i < len(s) && isPrint(s[i]) && s[i] != '"' {
		i++
	}
	if i >= len(s) || s[i] != '"' {
		return "", s, false
	}
	i++
	return s[:i], s[i:], true
}

func stringAppend(p []byte, s string) ([]byte, bool) {
	n, ok := stringCountEscapeChars(s)
	if !ok {
		return p, false
	}
	q := slices.Grow(p, len(`"`)+len(s)+n+len(`"`))
	q = append(q, '"')
	if n == 0 {
		q = append(q, s...)
		return append(q, '"'), true
	}
	q, ok = stringAppendEscape(q, s)
	if !ok {
		return p, false
	}
	return append(q, '"'), false
}

// stringCountEscapeChars return if s is a valid sf string value
// and how many characters that require escaping.
func stringCountEscapeChars(s string) (int, bool) {
	i, n := 0, 0
	for ; i < len(s) && isPrint(s[i]); i++ {
		if s[i] == '\\' || s[i] == '"' {
			n++
		}
	}
	return n, i >= len(s)
}

func StringValid(s string) bool {
	i := 0
	for i < len(s) && isPrint(s[i]) {
		i++
	}
	return i >= len(s)
}

func stringAppendEscape(p []byte, s string) ([]byte, bool) {
	q := p
	for i := 0; len(s) > 0; i = 1 {
		for i < len(s) && isPrint(s[i]) && s[i] != '"' && s[i] != '\\' {
			i++
		}
		if i >= len(s) {
			return append(q, s...), true
		}
		if !isPrint(s[i]) {
			return p, false
		}
		q, s = append(q, s[:i]...), s[i:]
		q = append(q, '\\')
	}
	return q, true
}

func stringAppendUnescape(p []byte, s string) ([]byte, bool) {
	q := p
	for i := 0; len(s) > 0; i = 1 {
		for i < len(s) && isPrint(s[i]) && s[i] != '\\' {
			i++
		}
		if i >= len(s) {
			return append(q, s...), true
		}
		if !isPrint(s[i]) {
			return p, false
		}
		// s[i] == '\\'
		q, s = append(q, s[:i]...), s[i:][1:]
		if len(s) <= 0 || (s[0] != '\\' && s[0] != '"') {
			return p, false
		}
	}
	return q, true
}
