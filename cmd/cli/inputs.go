package cli

import (
	"errors"
	"io"
	"os"
	"strings"
)

func NewInputReader(inputType string, args []string) (io.Reader, error) {
	switch inputType {
	case "console":
		return strings.NewReader(strings.Join(args, " ")), nil
	case "file":
		fileName := args[0]
		f, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		return nil, errors.New("unknown input type")
	}
}
