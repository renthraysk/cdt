package cdt

import "io"

type Reader struct {
	io.Reader
	io.Closer
}
