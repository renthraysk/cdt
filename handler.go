package cdt

import (
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

type CDT struct {
	current    atomic.Pointer[Resource]
	compendium *Compendium
}

func New(maxSize int, match, id string) (*CDT, error) {
	compendium, err := NewCompendium(maxSize, match, id)
	if err != nil {
		return nil, err
	}
	return &CDT{
		compendium: compendium,
	}, nil
}

func (cdt *CDT) Put(r io.Reader, lastModified time.Time) error {
	d, err := NewDictionary(r, lastModified)
	if err != nil {
		return err
	}
	cdt.compendium.Add(d)
	cdt.current.Store(&d.Resource)
	return nil
}

func (cdt *CDT) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet, http.MethodHead:
		res := cdt.current.Load()
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := cdt.compendium.Serve(w, r, res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodPut:
		if res := cdt.current.Load(); res != nil {
			if !res.evaluatePreconditions(r.Method, r.Header) {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
		}

		if err := cdt.Put(r.Body, time.Now()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
