package nio

import (
	"io"
)

func NewFnReader(
	fn func(bytes []byte) (int, error),
) io.Reader {
	return &fnReader{fn: fn}
}

type fnReader struct {
	fn func(bytes []byte) (int, error)
}

func (r fnReader) Read(p []byte) (int, error) {
	return r.fn(p)
}

func NewFnWriter(encode func(from []byte) error) io.Writer {
	return &fnWriter{fn: encode}
}

type fnWriter struct {
	fn func(from []byte) error
}

func (w fnWriter) Write(p []byte) (int, error) {
	err := w.fn(p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
