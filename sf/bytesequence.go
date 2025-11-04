package sf

import (
	"slices"

	"github.com/renthraysk/cdt/sf/b64"
)

func isBase64(c byte) bool {
	const base64 = digit | upper | lower | 1<<'+' | 1<<'/'

	return isASCII(c, base64%(1<<64), base64>>64)
}

// https://www.rfc-editor.org/rfc/rfc8941#name-byte-sequences

func ByteSequenceLen(n int) int {
	return len(":") + b64.EncodedLen(n) + len(":")
}

func ByteSequence(dst []byte, v []string) ([]byte, bool) {
	if len(v) != 1 {
		return nil, false
	}
	s, r, ok := byteSequenceCut(v[0])
	if !ok || len(s) < 2 || len(r) != 0 {
		return nil, false
	}
	return byteSequenceParse(dst, s)
}

func byteSequenceParse(dst []byte, s string) ([]byte, bool) {
	if len(s) < 2 || s[0] != ':' || s[len(s)-1] != ':' {
		return nil, false
	}
	n, err := b64.Decode(dst, s[1:len(s)-1])
	if err != nil {
		return nil, false
	}
	if !(len(dst) <= n && n <= cap(dst)) {
		return nil, false
	}
	return dst[:n], true
}

func ByteSequenceAppend(p, b []byte) []byte {
	p = slices.Grow(p, ByteSequenceLen(len(b)))
	p = append(p, ':')
	p = b64.AppendEncode(p, b)
	p = append(p, ':')
	return p
}

func byteSequenceCut(s string) (string, string, bool) {
	if len(s) <= 0 || s[0] != ':' {
		return "", s, false
	}
	i := 1
	for i < len(s) && isBase64(s[i]) {
		i++
	}
	for n := min(len(s), i+len("==")); i < n && s[i] == '='; {
		i++
	}
	if i >= len(s) || s[i] != ':' {
		return "", s, false
	}
	i++
	return s[:i], s[i:], true
}
