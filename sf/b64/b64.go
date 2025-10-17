package b64

import (
	"encoding/base64"
	"errors"
)

var (
	ErrPadding = errors.New("b64: incorrect padding")
	ErrLength  = errors.New("b64: invalid length")
)

func EncodedLen(n int) int {
	return base64.StdEncoding.EncodedLen(n)
}

// len(dst) is the minimum decoded length required
// cap(dst) is the maximum decoded length allowed
func Decode(dst []byte, s string) (int, error) {
	if len(s)%4 != 0 {
		return 0, ErrPadding
	}
	b := base64.StdEncoding.Strict()
	n := b.DecodedLen(len(s)) - paddingCount(s)
	if !(len(dst) <= n && n <= cap(dst)) {
		return 0, ErrLength
	}
	return b.Decode(dst[:cap(dst)], []byte(s))
}

func AppendEncode(dst []byte, src []byte) []byte {
	return base64.StdEncoding.AppendEncode(dst, src)
}

func paddingCount(s string) int {
	if len(s) <= 0 || s[len(s)-1] != '=' {
		return 0
	}
	if len(s) <= 1 || s[len(s)-2] != '=' {
		return 1
	}
	return 2
}
