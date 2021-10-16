package nio

import (
	"fmt"
	"os"
)

type FileCreator = func() (*os.File, error)
type fileWriter struct {
	fileName string
	file     *os.File
	onClose  func()
}

func NewFileWriter(fileName string, onClose func()) ClosableWriter {
	return &fileWriter{fileName: fileName, file: nil, onClose: onClose}
}
func (w *fileWriter) Write(p []byte) (n int, err error) {
	if w.file == nil {
		w.file, err = os.Create(w.fileName)
		if err != nil {
			return 0, fmt.Errorf("error creating file: %s", err)
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
