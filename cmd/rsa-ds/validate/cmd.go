package validate

import (
	"github.com/paulpaulych/crypto/internal/app/digital-sign"
)

type Cmd struct {
	PubKeyFile string `short:"p" long:"pub-key-file" description:"path to file containing public key" default:"rsa_pub.key"`
	Args       struct {
		SignedFile string `positional-arg-name:"signed-file" description:"path to file containing signed message"`
	} `positional-args:"true" required:"true"`
}

func (c *Cmd) Execute(_ []string) error {
	return digital_sign.Validate(c.Args.SignedFile, c.PubKeyFile)
}
