package nio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

)

const nameBufLen = 4


func WriteNamedStream(name string, content io.Reader, w io.Writer) error {
	nameSize:= len([]byte(name))
	nameSizeBuf := bytes.NewBuffer(make([]byte, nameBufLen))
	binary.BigEndian.PutUint32(nameSizeBuf.Bytes(), uint32(nameSize))	

	_, e := w.Write(nameSizeBuf.Bytes())
	if e != nil {
		return fmt.Errorf("can't write stream name size: %v", e)
	}

	_, e = io.Copy(w, strings.NewReader(name))
	if e != nil {
		return fmt.Errorf("can't write stream name: %v", e)
	}

	_, e = io.Copy(w, content)
	if e != nil {
		return fmt.Errorf("can't write stream content: %v", e)
	}

	return nil
}



type NamedStream struct {
	name string
	content io.Reader
}

func ReadNamedStream(r io.Reader) (*NamedStream, error) {
	nameSizeBuf := make([]byte, nameBufLen)

	_, e := r.Read(nameSizeBuf)
	if e != nil {
		return nil, fmt.Errorf("can't read file name size: %v", e)
	}

	nameSize := binary.BigEndian.Uint32(nameSizeBuf)	
	
	nameBuf := make([]byte, nameSize)
	
	_, e = r.Read(nameBuf)
	if e != nil {
		return nil, fmt.Errorf("can't read file name: %v", e)
	}

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		io.Copy(pw, r)
	}()

	return &NamedStream{
		name: string(nameBuf),
		content: pr,
	}, nil
}
