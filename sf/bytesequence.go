package sf

import (
	"slices"

	"github.com/renthraysk/cdt/sf/b64"
)

// https://www.rfc-editor.org/rfc/rfc8941#name-byte-sequences

func ByteSequenceLen(n int) int {
	return len(":") + b64.EncodedLen(n) + len(":")
}

func ByteSequence(dst []byte, v []string) ([]byte, bool) {
	if len(v) != 1 {
		return nil, false
	}
	s := v[0]
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

func AppendByteSequence(p, b []byte) []byte {
	p = slices.Grow(p, ByteSequenceLen(len(b)))
	p = append(p, ':')
	p = b64.AppendEncode(p, b)
	p = append(p, ':')
	return p
}
