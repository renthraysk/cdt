package sf

import (
	"math"
	"strconv"
)

const (
	integerDigits       = 15
	integerMin    int64 = 1 - 1e15
	integerMax    int64 = 1e15 - 1

	decimalDigits             = 12
	decimalFracDigits         = 3
	decimalRound              = 1e-3
	decimalMin        float64 = decimalRound - 1e12
	decimalMax        float64 = 1e12 - decimalRound
)

func numericCut(s string) (string, string, bool) {
	if len(s) <= 0 {
		return "", s, false
	}
	i := 0
	if s[0] == '-' {
		i = len("-")
	}
	if i >= len(s) || !isDigit(s[i]) {
		return "", s, false
	}
	n := i
	i++
	for i < len(s) && isDigit(s[i]) {
		i++
	}
	if i >= len(s) || s[i] != '.' {
		// integer
		n += integerDigits
		if i > n {
			return "", s, false
		}
		if i >= len(s) {
			return s, "", true
		}
		return s[:i], s[i:], true
	}
	// s[i] == '.'
	// decimal
	n += decimalDigits
	if i > n {
		return "", s, false
	}
	i++
	for n := min(len(s), i+decimalFracDigits); i < n && isDigit(s[i]); {
		i++
	}
	if i >= len(s) {
		// i+decimalFracDigits was < len(s)
		return s, "", true
	}
	if isDigit(s[i]) {
		return "", s, false
	}
	return s[:i], s[i:], true
}

func integerParse(s string) (int64, bool) {
	i, err := strconv.ParseInt(s, 10, 64)
	return i, err == nil
}

func decimalParse(s string) (float64, bool) {
	f, err := strconv.ParseFloat(s, 64)
	return f, err == nil
}

func integerAppend(p []byte, i int64) ([]byte, bool) {
	if integerMin <= i && i <= integerMax {
		return strconv.AppendInt(p, i, 10), true
	}
	return p, false
}

func decimalAppend(p []byte, f float64) ([]byte, bool) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return p, false
	}
	r := math.RoundToEven(f/decimalRound) * decimalRound
	if decimalMin <= r && r <= decimalMax {
		return strconv.AppendFloat(p, r, 'f', decimalFracDigits, 64), true
	}
	return p, false
}
