package nio

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type FileData struct {
	Name string
	Content io.Reader
}

func NewFileData(f *os.File) *FileData {
	return &FileData{
		Name: f.Name(),
		Content: f,
	}
} 

func (s *FileData) WriteTo(w io.Writer) (int64, error) {
	nameSize:= len([]byte(s.Name))
	nameSizeBuf := []byte{byte(nameSize)}
	totalWrote := int64(0)
	n1, e := w.Write(nameSizeBuf)
	if e != nil {
		return 0, fmt.Errorf("can't write stream name size: %v", e)
	}
	totalWrote += int64(n1)

	n2, e := io.Copy(w, strings.NewReader(s.Name))
	if e != nil {
		return 0, fmt.Errorf("can't write stream name: %v", e)
	}
	totalWrote += int64(n2)

	n2, e = io.Copy(w, s.Content)
	if e != nil {
		return 0, fmt.Errorf("can't write stream content: %v", e)
	}

	return totalWrote + int64(n2), nil
}

func ReadFileData(r io.Reader) (*FileData, error) {
	nameSizeBuf := make([]byte, 1)

	_, e := r.Read(nameSizeBuf)
	if e != nil {
		return nil, fmt.Errorf("can't read file name size: %v", e)
	}
	
	nameSize := int(nameSizeBuf[0])
	nameBuf := make([]byte, nameSize)

	fmt.Println("name size", nameSize)

	for i := 0; i < nameSize; {	
		n, e := r.Read(nameBuf[i:])
		if e != nil {
			return nil, fmt.Errorf("can't read file name: %v", e)
		}
		i += n
	}

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		io.Copy(pw, r)
	}()

	return &FileData{
		Name: string(nameBuf),
		Content: pr,
	}, nil
}

func (s *FileData) WriteToFile() error {
	f, e := os.Create(s.Name + ".received")
	if e != nil {
		return e
	}
	_, e = io.Copy(f, s.Content)
	if e != nil {
		return e
	}
	return nil
}