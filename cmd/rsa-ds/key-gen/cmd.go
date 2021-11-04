package key_gen

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app"
)

type Cmd struct {
	P *cli.BigIntOpt `short:"P" description:"large prime number" required:"true"`
	Q *cli.BigIntOpt `short:"Q" description:"large prime number" required:"true"`
}

func (c *Cmd) Execute(_ []string) error {
	return app.RsaKeyGen(c.P.Value, c.Q.Value)
}
