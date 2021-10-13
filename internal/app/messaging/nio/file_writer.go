package nio

import (
	"errors"
	"fmt"
	"os"
)

type FileCreator = func() (*os.File, error)
type fileWriter struct {
	newFile FileCreator
	file    *os.File
	onClose func()
}

func NewFileWriter(newFile FileCreator, onClose func()) ClosableWriter {
	return &fileWriter{newFile: newFile, file: nil, onClose: onClose}
}
func (w *fileWriter) Write(p []byte) (n int, err error) {
	if w.file == nil {
		w.file, err = w.newFile()
		if err != nil {
			errMsg := fmt.Sprintf("error creating file: %s", err)
			return 0, errors.New(errMsg)
		}
	}
	return w.file.Write(p)
}
func (w *fileWriter) Close() error {
	if w.file == nil {
		return nil
	}
	err := w.file.Close()
	if err != nil {
		return err
	}
	w.onClose()
	return nil
}
