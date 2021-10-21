package sign

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app/digital-sign"
)

type Cmd struct {
	SecretFile string `long:"secret" description:"path to file containing secret key" default:"rsa.key"`
	Input      string `short:"i" long:"input" choice:"file" choice:"console" description:"input type: console or file" default:"console"`
	Args       struct {
		Msg string `positional-arg-name:"message" description:"message text or name of file when -i=file specified"`
	} `positional-args:"true" required:"true"`
}

func (c *Cmd) Execute(_ []string) error {
	msgReader, e := cli.NewInputReader(c.Input, c.Args.Msg)
	if e != nil {
		return e
	}
	return digital_sign.Sign(msgReader, c.SecretFile)
}
