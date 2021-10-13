package nio

import (
	"fmt"
	"strings"
)

type consoleWriter struct {
	buf *strings.Builder
}

func NewConsoleWriter(prefix []byte) ClosableWriter {
	buf := new(strings.Builder)
	buf.Write(prefix)
	return &consoleWriter{buf: buf}
}
func (w *consoleWriter) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}
func (w *consoleWriter) Close() error {
	fmt.Println(w.buf.String())
	return nil
}
