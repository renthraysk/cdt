package cdt

import (
	"errors"
	"maps"
	"net/http"
	"sync"

	"github.com/renthraysk/encoding"
)

// Compendium represents a ring buffer of compression dictionaries with thread-safe access.
type Compendium struct {
	id              string
	useAsDictionary string // Use-As-Dictionary header value

	mu       sync.RWMutex
	ring     []*Dictionary
	addPoint int
}

// NewCompendium creates a new Compendium instance with the specified maximum number of
// Dictionaries available for compression.
// match, and id refer to the Use-As-Dictionary header values.
// All dictionaries in a compendium will share the same id, with the
// Available-Dictionary used to pick an individual dictionary if
// Dictionary-ID matches the Compendium's id.
func NewCompendium(maxSize int, match, id string) (*Compendium, error) {
	useAsDictionary, err := UseAsDictionary(match, id)
	if err != nil {
		return nil, err
	}
	return &Compendium{
		id:              id,
		useAsDictionary: useAsDictionary,
		ring:            make([]*Dictionary, 0, maxSize),
	}, nil
}

func (c *Compendium) Serve(w http.ResponseWriter, r *http.Request, res *Resource) error {
	if r.Method != http.MethodHead && r.Method != http.MethodGet {
		return errors.New("Serve called with unexpected method")
	}

	maps.Insert(w.Header(), res.Cached.headers)
	w.Header().Set("Vary", "accept-encoding, available-dictionary")

	if !res.evaluatePreconditions(r.Method, r.Header) {
		w.WriteHeader(http.StatusNotModified)
		return nil
	}

	maps.Insert(w.Header(), res.headers)
	w.Header().Set("Use-As-Dictionary", c.useAsDictionary)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	ae := encoding.ParseAcceptEncoding(r.Header.Get("Accept-Encoding"))
	if ae.Contains(encoding.DCZ) {
		if ad, ok := availableDictionary(r.Header); ok {
			id, _ := dictionaryID(r.Header)
			if d := c.Get(ad, id); d != nil {
				w.Header().Set("Content-Encoding", "dcz")
				return d.Encode(w, res)
			}
		}
	}
	switch {
	case ae.Contains(encoding.Zstd):
		w.Header().Set("Content-Encoding", "zstd")
		return Zstd(w, res)

	case ae.Contains(encoding.Gzip):
		w.Header().Set("Content-Encoding", "gzip")
		return Gzip(w, res)
	}
	w.Header().Del("Content-Encoding")
	_, err := res.WriteTo(w)
	return err
}

func (c *Compendium) Add(newD *Dictionary) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cap(c.ring) <= 0 {
		return
	}
	for _, d := range c.ring {
		if d.SHA256() == newD.SHA256() {
			return
		}
	}
	if len(c.ring) <= 0 || len(c.ring) < cap(c.ring) {
		c.ring = append(c.ring, newD)
		return
	}
	i := c.addPoint
	if i >= len(c.ring) {
		i = 0
	}
	c.ring[i] = newD
	i++
	if i >= len(c.ring) {
		i = 0
	}
	c.addPoint = i
}

func (c *Compendium) Get(ad []byte, id string) *Dictionary {
	if c.id != id {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, d := range c.ring {
		if d.SHA256() == string(ad) {
			return d
		}
	}
	return nil
}
