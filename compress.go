package cdt

import (
	"io"
	"sync"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
)

var gzipPool = sync.Pool{New: func() any { return gzip.NewWriter(nil) }}

func Gzip(w io.Writer, r io.WriterTo) error {
	z := gzipPool.Get().(*gzip.Writer)
	z.Reset(w)
	defer gzipPool.Put(z)
	defer z.Reset(nil)
	defer z.Close()
	if _, err := r.WriteTo(z); err != nil {
		return err
	}
	return z.Close()
}

var zstdPool = sync.Pool{New: func() any {
	z, err := zstd.NewWriter(nil)
	if err != nil {
		return nil
	}
	return z
}}

func Zstd(w io.Writer, r io.WriterTo) error {
	z := zstdPool.Get().(*zstd.Encoder)
	if z == nil {
		panic("no zstd encoder")
	}
	z.Reset(w)
	defer zstdPool.Put(z)
	defer z.Reset(nil)
	defer z.Close()
	if _, err := r.WriteTo(z); err != nil {
		return err
	}
	return z.Close()
}
