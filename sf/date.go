package sf

import "time"

func dateCut(s string) (string, string, bool) {
	if len(s) < len("@0") || s[0] != '@' || !isDigit(s[1]) {
		return "", s, false
	}
	i := 2
	for i < len(s) && isDigit(s[i]) {
		i++
	}
	return s[:i], s[i:], true
}

func dateParse(s string) (time.Time, bool) {
	if len(s) >= len("@0") && s[0] == '@' && isDigit(s[1]) {
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
