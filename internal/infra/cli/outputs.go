package cli

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/nio"
	"log"
	"net"
	"time"
)

type OutputFactory = func(from net.Addr) nio.ClosableWriter

func NewOutputFactory(outputType string) (OutputFactory, error) {
	switch outputType {
	case "console":
		return func(from net.Addr) nio.ClosableWriter {
			prefix := fmt.Sprintln("RECEIVED MESSAGE FROM", from)
			return nio.NewConsoleWriter([]byte(prefix))
		}, nil
	case "file":
		return func(from net.Addr) nio.ClosableWriter {
			fName := fmt.Sprintf("msg_from_%s_at_%v.txt", from.String(), time.Now().UnixMilli())
			return nio.NewFileWriter(fName, func() {
				log.Printf("received message saved to file %s", fName)
			})
		}, nil
	default:
		return nil, errors.New("invalid outputType type")
	}
}
