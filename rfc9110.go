package cdt

import (
	"io"
	"net/http"
	"strings"
	"time"
)

type Cached struct {
	lastModifiedTime time.Time
	lastModified     []string
	cacheControl     []string
	eTag             []string
}

type Resource struct {
	Cached
	contentType []string
	data        []byte
}

func (r *Resource) headers(yield func(string, []string) bool) {
	yield("Content-Type", r.contentType)
}

func (r *Resource) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(r.data)
	return int64(n), err
}

func (c *Cached) headers(yield func(string, []string) bool) {
	if len(c.cacheControl) > 0 && !yield("Cache-Control", c.cacheControl) {
		return
	}
	if len(c.lastModified) > 0 && !yield("Last-Modified", c.lastModified) {
		return
	}
	if len(c.eTag) > 0 {
		yield("Etag", c.eTag)
	}
}

// evaluatePreconditions implements
// https://www.rfc-editor.org/rfc/rfc9110.html#name-preconditions
func (c *Cached) evaluatePreconditions(method string, r http.Header) bool {

	// https://www.rfc-editor.org/rfc/rfc9110.html#name-precedence-of-preconditions
	// Step 1
	if m, ok := c.ifMatch(r); ok {
		if !m {
			return false
		}
	} else /* 2 */ if us, ok := c.ifUnmodifiedSince(r); ok && !us {
		return false
	}

	switch method {
	case http.MethodGet, http.MethodHead:
		// Step 3
		if nm, ok := c.ifNoneMatch(r, true); ok {
			if !nm {
				return false // (304 Not Modified)
			}
		} else /* 4 */ if ms, ok := c.ifModifiedSince(r); ok && !ms {
			return false // (304 Not Modified)
		}
		/* Step 5
		if method != http.MethodGet {
			break
		}
		if rg, ok := Range(r); ok {
			if ir, ok := c.ifRange(r); ok && ir {
				break
			}
			r.Del("Range")
		}
		*/
	default:
		// Step 3
		if nm, ok := c.ifNoneMatch(r, false); ok && !nm {
			return false // (412 Precondition failed)
		}
	}
	// Step 6
	return true
}

func eTagNormalize(s string) (string, bool) {
	const prefix = `W/`

	s = trimOWS(s)

	weak := false
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		weak = true
		s = s[len(prefix):]
	}
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s, weak
}

type Etags string

func (e Etags) Tags(yield func(eTag string, weak bool) bool) {
	list := string(e)
	for i := strings.IndexByte(list, ','); i >= 0; i = strings.IndexByte(list, ',') {
		if s, weak := eTagNormalize(list[:i]); len(s) > 0 && !yield(s, weak) {
			return
		}
		list = list[i:][len(","):]
	}
	if s, weak := eTagNormalize(list); len(s) > 0 {
		yield(s, weak)
	}
}

// https://www.rfc-editor.org/rfc/rfc9110.html#name-if-match
func (c *Cached) ifMatch(r http.Header) (eval bool, present bool) {
	switch text := r.Get("If-Match"); text {
	case "":
		return false, false
	case "*":
		return true, true
	default:
		if len(c.eTag) == 1 {
			if eTag, weak := eTagNormalize(c.eTag[0]); !weak {
				for e, w := range Etags(text).Tags {
					if !w && e == eTag {
						return true, true
					}
				}
			}
		}
	}
	return false, true
}

// https://www.rfc-editor.org/rfc/rfc9110.html#name-if-none-match
func (c *Cached) ifNoneMatch(r http.Header, isSafeMethod bool) (eval bool, present bool) {
	switch text := r.Get("If-None-Match"); text {
	case "":
		return false, false
	case "*":
		return false, true
	default:
		if len(c.eTag) == 1 {
			if eTag, weak := eTagNormalize(c.eTag[0]); !isSafeMethod || !weak {
				for e, w := range Etags(text).Tags {
					if (!isSafeMethod || !w) && e == eTag {
						return false, true
					}
				}
			}
		}
	}
	return true, true
}

func (c *Cached) ifUnmodifiedSince(r http.Header) (eval bool, present bool) {
	if text := r.Get("If-Unmodified-Since"); text != "" {
		if t, err := http.ParseTime(text); err == nil {
			return !c.lastModifiedTime.After(t), true
		}
	}
	return false, false
}

func (c *Cached) ifModifiedSince(r http.Header) (eval bool, present bool) {
	if text := r.Get("If-Modified-Since"); text != "" {
		if t, err := http.ParseTime(text); err == nil {
			return c.lastModifiedTime.After(t), true
		}
	}
	return false, false
}

func (c *Cached) ifRange(r http.Header) (eval bool, present bool) {
	if text := r.Get("If-Range"); text != "" {
		if len(c.eTag) == 1 {
			tag, weak := eTagNormalize(c.eTag[0])
			if t, w := eTagNormalize(text); weak == w && tag == t {
				return true, true
			}
		}
		if !c.lastModifiedTime.IsZero() {
			if t, err := http.ParseTime(text); err == nil {
				return t.Unix() == c.lastModifiedTime.Unix(), true
			}
		}
	}
	return false, false
}

func isOWS(c byte) bool { return c == ' ' || c == '\t' }

func trimOWS(s string) string {
	i := 0
	for i < len(s) && isOWS(s[i]) {
		i++
	}
	if i >= len(s) {
		return ""
	}
	n := len(s) - 1
	for n > i && isOWS(s[n]) {
		n--
	}
	n++
	return s[i:n]
}
