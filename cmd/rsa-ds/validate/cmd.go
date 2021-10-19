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
	flagsSpec := cli.NewFlagSpec(conf.CmdName(), map[string]string{
		"pub": "path to file containing public key",
	})

	flags, err := flagsSpec.Parse(args)
	if err != nil {
		return nil, err
	}

	if len(flags.Args) < 1 {
		return nil, cli.NewCmdConfError("args required: [signed file]", nil)
	}

	pubFName := flags.Flags["pub"].Get()

	if pubFName == nil {
		return nil, cli.NewCmdConfError("required flag: -pub", nil)
	}
	return &Cmd{signedFile: flags.Args[0], pubKeyFile: *pubFName}, nil
}

type Cmd struct {
	signedFile string
	pubKeyFile string
}

func (cmd *Cmd) Run() error {
	return digital_sign.Validate(cmd.signedFile, cmd.pubKeyFile)
}
