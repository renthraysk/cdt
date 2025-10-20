package cdt

import (
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

// SelfPack: An HTTP resource that serves as its own compression dictionary.
// When a client request supports 'dcz' Content-Encoding and includes an
// Available-Dictionary matching a retained prior version, the latest resource
// version is compressed using that prior version as the compression dictionary.
type SelfPack struct {
	current    atomic.Pointer[Resource]
	compendium *Compendium
}

func NewSelfPack(maxSize int, match, id string) (*SelfPack, error) {
	compendium, err := NewCompendium(maxSize, match, id)
	if err != nil {
		return nil, err
	}
	return &SelfPack{
		compendium: compendium,
	}, nil
}

func (sp *SelfPack) Put(r io.Reader, lastModified time.Time) error {
	d, err := NewDictionary(r, lastModified)
	if err != nil {
		return err
	}
	sp.compendium.Add(d)
	sp.current.Store(&d.Resource)
	return nil
}

func (sp *SelfPack) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet, http.MethodHead:
		res := sp.current.Load()
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := sp.compendium.Serve(w, r, res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodPut:
		if res := sp.current.Load(); res != nil {
			if !res.evaluatePreconditions(r.Method, r.Header) {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
		}

		if err := sp.Put(r.Body, time.Now()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
