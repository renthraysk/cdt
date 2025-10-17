package cdt

import (
	"crypto/sha256"
	"errors"
	"net/http"

	"github.com/renthraysk/cdt/sf"
)

const DictionaryIDMaxLength int = 1024

var (
	ErrMatchRequired = errors.New("Use-As-Dictionary: match required")
	ErrMatchInvalid  = errors.New("Use-As-Dictionary: match invalid")
	ErrIDInvalid     = errors.New("Use-As-Dictionary: id invalid")
)

func dictionaryID(r http.Header) (string, bool) {
	id, ok := sf.String(r["Dictionary-ID"])
	if !ok || len(id) > DictionaryIDMaxLength {
		return "", false
	}
	return id, true
}

func availableDictionary(r http.Header) ([]byte, bool) {
	dst := make([]byte, sha256.Size, sha256.Size)
	return sf.ByteSequence(dst, r["Available-Dictionary"])
}

// https://www.ietf.org/archive/id/draft-ietf-httpbis-compression-dictionary-19.html#name-use-as-dictionary
// type is raw
func appendUseAsDictionaryRaw(p []byte, match, id string) ([]byte, error) {

	// https://www.ietf.org/archive/id/draft-ietf-httpbis-compression-dictionary-19.html#name-match
	if len(match) == 0 {
		return p, ErrMatchRequired
	}
	q, ok := sf.AppendKeyString(p, "match", match)
	if !ok {
		return p, ErrMatchInvalid
	}

	// https://www.ietf.org/archive/id/draft-ietf-httpbis-compression-dictionary-19.html#name-id
	if len(id) > 0 {
		if len(id) > DictionaryIDMaxLength {
			return p, ErrIDInvalid
		}
		q = append(q, ' ')
		if q, ok = sf.AppendKeyString(q, "id", id); !ok {
			return p, ErrIDInvalid
		}
	}
	return q, nil
}

func UseAsDictionary(match, id string) (string, error) {
	b, err := appendUseAsDictionaryRaw(make([]byte, 0, 64), match, id)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
