package key_gen

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app/digital-sign"
)

type Cmd struct {
	P *cli.BigIntOpt `short:"p" description:"large prime number" required:"true"`
	Q *cli.BigIntOpt `short:"q" description:"large prime number" required:"true"`
}

func (c *Cmd) Execute(_ []string) error {
	return digital_sign.GenerateKeys(c.P.Value, c.Q.Value)
}
