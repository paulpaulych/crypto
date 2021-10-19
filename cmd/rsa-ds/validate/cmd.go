package validate

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	digital_sign "github.com/paulpaulych/crypto/internal/app/digital-sign"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "validate"
}

func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	var opts struct {
		PubKeyFile string `short:"p" long:"pub-key-file" description:"path to file containing public key" default:"rsa_pub.key"`
		Args       struct {
			SignedFile string `positional-arg-name:"signed-file" description:"path to file containing signed message"`
		} `positional-args:"true" required:"true"`
	}
	_, err := cli.ParseFlagsOfCmd(conf.CmdName(), &opts, args)
	if err != nil {
		return nil, err
	}

	return &Cmd{signedFile: opts.Args.SignedFile, pubKeyFile: opts.PubKeyFile}, nil
}

type Cmd struct {
	signedFile string
	pubKeyFile string
}

func (cmd *Cmd) Run() error {
	return digital_sign.Validate(cmd.signedFile, cmd.pubKeyFile)
}
