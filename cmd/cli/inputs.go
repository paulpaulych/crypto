package cli

import (
	"errors"
	"io"
	"os"
	"strings"
)

func NewInputReader(inputType string, msg string) (io.Reader, error) {
	switch inputType {
	case "console":
		return strings.NewReader(msg), nil
	case "file":
		fileName := msg
		f, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		return nil, errors.New("unknown input type")
	}
}
