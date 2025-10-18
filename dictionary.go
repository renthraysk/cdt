package cdt

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/renthraysk/cdt/sf"
)

const dbzmagic = "\x5e\x2a\x4d\x18\x20\x00\x00\x00"

type Dictionary struct {
	Resource
	header []byte
	pool   sync.Pool
}

func NewDictionary(r io.Reader, lastModified time.Time) (*Dictionary, error) {
	var b bytes.Buffer

	if _, err := b.ReadFrom(r); err != nil {
		return nil, err
	}
	data := b.Bytes()

	// attempt to create a writer to ensure zstd is happy with the dict
	z, err := zstd.NewWriter(nil, zstd.WithEncoderDictRaw(0, data))
	if err != nil {
		return nil, fmt.Errorf("NewDictionary: %v", err)
	}

	lm := lastModified.UTC().Truncate(time.Second)

	d := &Dictionary{
		Resource: Resource{
			contentType: []string{"text/plain"}, // @TODO configurable
			Cached: Cached{
				lastModifiedTime: lm,
				lastModified:     []string{lm.Format(http.TimeFormat)},
				eTag:             nil,
				cacheControl:     []string{"max-age=600, stale-whilst-revalidate=300"}, // @TODO configurable
			},
			data: data,
		},
		header: makeHeader(data),
		pool: sync.Pool{New: func() any {
			z, err := zstd.NewWriter(nil, zstd.WithEncoderDictRaw(0, data))
			if err != nil {
				panic(err.Error()) // @TODO eliminate
			}
			return z
		}},
	}
	d.pool.Put(z)
	return d, nil
}

func (d *Dictionary) sha256() []byte {
	return d.header[len(dbzmagic):][:sha256.Size]
}

func (d *Dictionary) SHA256() string {
	return string(d.sha256())
}

func (d *Dictionary) AvailableDictionary() string {
	p := make([]byte, 0, sf.ByteSequenceLen(sha256.Size))
	return string(sf.AppendByteSequence(p, d.sha256()))
}

func (d *Dictionary) Encode(w io.Writer, r io.WriterTo) error {
	if _, err := w.Write(d.header); err != nil {
		return err
	}
	z := d.pool.Get().(*zstd.Encoder)
	if z == nil {
		panic("no zstd encoder")
	}
	z.Reset(w)
	defer d.pool.Put(z)
	defer z.Reset(nil)
	defer z.Close()
	if _, err := r.WriteTo(z); err != nil {
		return err
	}
	return z.Close()
}

func makeHeader(data []byte) []byte {
	p := make([]byte, len(dbzmagic), len(dbzmagic)+sha256.Size)
	copy(p, dbzmagic)
	h := sha256.New()
	h.Write(data)
	return h.Sum(p)
}
