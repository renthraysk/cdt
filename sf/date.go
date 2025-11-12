package sf

import "time"

func dateCut(s string) (string, string, bool) {
	if len(s) < len("@0") || s[0] != '@' {
		return "", s, false
	}
	i := len("@")
	if s[1] == '-' {
		i = len("@-")
	}
	if i >= len(s) || !isDigit(s[i]) {
		return "", s, false
	}
	for n := min(len(s), i+integerDigits); i < n && isDigit(s[i]); {
		i++
	}
	if i >= len(s) {
		return s, "", true
	}
	if isDigit(s[i]) || s[i] == '.' {
		return "", s, false
	}
	return s[:i], s[i:], true
}

func dateParse(s string) (time.Time, bool) {
	if len(s) >= len("@0") && s[0] == '@' {
		if s, ok := integerParse(s[1:]); ok {
			return time.Unix(s, 0), true
		}
	}
	return time.Time{}, false
}

func dateAppend(p []byte, t time.Time) ([]byte, bool) {
	q := append(p, '@')
	if q, ok := integerAppend(q, t.Unix()); ok {
		return q, true
	}
	return p, false
}
