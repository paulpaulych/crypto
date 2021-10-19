package sign

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app/digital-sign"
	"io"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "sign"
}

func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	flagsSpec := cli.NewFlagSpec(conf.CmdName(), map[string]string{
		"secret": "path to file containing secret key",
		"i":      "message input type: file or arg",
	})

	flags, err := flagsSpec.Parse(args)
	if err != nil {
		return nil, err
	}

	if len(flags.Args) < 1 {
		return nil, cli.NewCmdConfError("args required: [message]", nil)
	}

	secretFName := flags.Flags["secret"].Get()
	input := flags.Flags["i"].GetOr("console")

	msgReader, e := cli.NewInputReader(input, flags.Args)
	if e != nil {
		return nil, cli.NewCmdConfError(e.Error(), nil)
	}

	if secretFName == nil {
		return nil, cli.NewCmdConfError("required flag: -secret", nil)
	}
	return &Cmd{msg: msgReader, secretFile: *secretFName}, nil
}

type Cmd struct {
	secretFile string
	msg        io.Reader
}

func (cmd *Cmd) Run() error {
	return digital_sign.Sign(cmd.msg, cmd.secretFile)
}
