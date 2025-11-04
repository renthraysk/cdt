package sf

import (
	"math/bits"
	"strings"
	"unicode/utf8"
)

func displayStringCut(s string) (string, string, bool) {
	const prefix = `%"`

	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		n, ok := 0, true
		for i := len(prefix); ok && i < len(s); i += n {
			for i < len(s) && isPrint(s[i]) && s[i] != '"' && s[i] != '%' {
				i++
			}
			if i >= len(s) || !isPrint(s[i]) {
				break
			}
			if s[i] == '"' {
				i++
				return s[:i], s[i:], true
			}
			// s[i] == '%'
			_, n, ok = percentDecode(s[i:])
		}
	}
	return "", s, false
}

func DisplayStringParse(s string) (string, bool) {
	const prefix = `%"`

	if len(s) < len(prefix)+len(`"`) || s[:len(prefix)] != prefix || s[len(s)-1] != '"' {
		return "", false
	}
	return displayStringDecode(s[len(prefix) : len(s)-1])
}

func displayStringDecode(s string) (string, bool) {
	i := 0
	for i < len(s) && isPrint(s[i]) && s[i] != '%' {
		i++
	}
	if i >= len(s) {
		return s, true
	}
	if !isPrint(s[i]) {
		return "", false
	}
	// s[i] == '%'
	r, n, ok := percentDecode(s[i:])
	if !ok {
		return "", false
	}

	var b strings.Builder

	b.Grow(len(s))
	b.WriteString(s[:i])
	b.WriteRune(r)
	s = s[i:][n:]
	for len(s) > 0 {
		i := 0
		for i < len(s) && isPrint(s[i]) && s[i] != '%' {
			i++
		}
		if i >= len(s) {
			b.WriteString(s)
			return b.String(), true
		}
		if !isPrint(s[i]) {
			return "", false
		}
		b.WriteString(s[:i])
		s = s[i:]
		r, n, ok := percentDecode(s)
		if !ok {
			return "", false
		}
		b.WriteRune(r)
		s = s[n:]
	}
	return b.String(), true
}

func appendDisplayStringEncode(p []byte, s string) ([]byte, bool) {
	const h = "0123456789abcdef"

	b := make([]byte, 0, utf8.UTFMax)
	q := p
	for len(s) > 0 {
		i := 0
		for i < len(s) && isPrint(s[i]) && s[i] != '%' && s[i] != '"' {
			i++
		}
		if i >= len(s) {
			return append(q, s...), true
		}
		r, n := utf8.DecodeRuneInString(s[i:])
		if n <= 1 && r == utf8.RuneError {
			return p, false
		}
		q, s = append(q, s[:i]...), s[i:]
		for _, x := range utf8.AppendRune(b, r) {
			q = append(q, '%', h[x>>4], h[x&0xF])
		}
		s = s[n:]
	}
	return q, true
}

// percentLen parses a percent-encoded UTF-8 sequence from the start of s.
// It returns the length of the valid sequence and true if successful, or 0 and false otherwise.
func percentDecode(s string) (rune, int, bool) {
	if len(s) < 3 || s[0] != '%' {
		return utf8.RuneError, 0, false
	}
	x, ok := hexByteDecode(s[1], s[2])
	if !ok {
		return utf8.RuneError, 0, false
	}
	if x <= 0x7F {
		return rune(x), 3, true // x = 0b0xxxxxxx, ASCII
	}
	n := bits.LeadingZeros8(^x)
	if n < 2 || n > utf8.UTFMax {
		return utf8.RuneError, 0, false
	}
	// x = 0b110xxxxx, 0b1110xxxx, 0b11110xxx
	// 2, 3, or 4 byte UTF-8 sequence
	encLen := n * len("%xx")
	if len(s) < encLen {
		return utf8.RuneError, 0, false
	}
	b := make([]byte, 1, utf8.UTFMax)
	b[0] = x
	for y := range percentEncoded(s[3:encLen]).Decode {
		b = append(b, y)
	}
	if len(b) == n {
		if r, nn := utf8.DecodeRune(b); n == nn {
			return r, encLen, true
		}
	}
	return utf8.RuneError, 0, false
}

type percentEncoded string

func (s percentEncoded) Decode(yield func(byte) bool) {
	for len(s) >= 3 && s[0] == '%' {
		x, ok := hexByteDecode(s[1], s[2])
		s = s[3:]
		if !ok || !yield(x) {
			break
		}
	}
}

func hexByteDecode(a, b uint8) (uint8, bool) {
	if x, ok := hexDigitDecode(a); ok {
		y, ok := hexDigitDecode(b)
		return x<<4 | y, ok
	}
	return 0, false
}

func hexDigitDecode(c uint8) (uint8, bool) {
	if '0' <= c && c <= '9' {
		return c - '0', true
	}
	if 'a' <= c && c <= 'f' {
		return c - 'a' + 10, true
	}
	return c - 'A' + 10, 'A' <= c && c <= 'F'
}
