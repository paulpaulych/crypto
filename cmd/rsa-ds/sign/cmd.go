package sign

import "github.com/paulpaulych/crypto/internal/app"

type Cmd struct {
	SecretFile string `short:"s" long:"secret" description:"path to file containing secret key" required:"true"`
	Args       struct {
		File string `positional-arg-name:"file" description:"path to file"`
	} `positional-args:"true" required:"true"`
}

func (c *Cmd) Execute(_ []string) error {
	return app.RsaSign(c.Args.File, c.SecretFile)
}
