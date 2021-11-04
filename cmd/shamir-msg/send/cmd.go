package send

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app"
)

type Cmd struct {
	P     *cli.BigIntOpt `short:"P" description:"large prime number" required:"true"`
	Args  struct {
		Addr string `positional-arg-name:"target" description:"target host:port. Example: 127.0.0.1:1234"`
		File  string `positional-arg-name:"file" description:"name of file"`
	} `positional-args:"true" required:"true"`
}

func (c *Cmd) Execute(_ []string) error {
	return app.ShamirSend(c.Args.Addr, c.Args.File, c.P.Value)
}
