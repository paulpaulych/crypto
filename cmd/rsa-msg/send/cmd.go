package send

import (
	"github.com/paulpaulych/crypto/internal/app"
)

type Cmd struct {
	BobPub string `short:"b" long:"bob-pub" description:"path to file containing destination public key" required:"true"`
	Args   struct {
		Addr string `positional-arg-name:"target address" description:"target host:port. Example: 127.0.0.1:1234"`
		File  string `positional-arg-name:"file" description:"path to file"`
	} `positional-args:"true" required:"true"`
}

func (c *Cmd) Execute(_ []string) error {
	return app.RsaSend(c.Args.Addr, c.Args.File, c.BobPub)
}
