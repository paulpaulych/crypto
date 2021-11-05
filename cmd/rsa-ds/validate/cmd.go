package validate

import "github.com/paulpaulych/crypto/internal/app"

type Cmd struct {
	PubKeyFile string `short:"p" long:"pub-key-file" description:"path to file containing public key" required:"true"`
	SignatureFile string `short:"s" long:"signature-file" description:"path to file containing signature" required:"true"`
	Args       struct {
		MsgFile string `positional-arg-name:"msg-file" description:"path to file containing original message"`
	} `positional-args:"true" required:"true"`
}

func (c *Cmd) Execute(_ []string) error {
	return app.RsaValidate(c.Args.MsgFile, c.SignatureFile, c.PubKeyFile)
}
