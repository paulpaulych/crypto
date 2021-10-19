package sign

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	digital_sign "github.com/paulpaulych/crypto/internal/app/digital-sign"
	"io"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "sign"
}

func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	var opts struct {
		SecretFile string `long:"secret" description:"path to file containing secret key" default:"rsa.key"`
		Input      string `short:"i" long:"input" choice:"file" choice:"console" description:"input type: console or file" default:"console"`
		Args       struct {
			Msg string `positional-arg-name:"message" description:"message text or name of file when -i=file specified"`
		} `positional-args:"true" required:"true"`
	}

	_, err := cli.ParseFlagsOfCmd(conf.CmdName(), &opts, args)
	if err != nil {
		return nil, err
	}

	msgReader, e := cli.NewInputReader(opts.Input, opts.Args.Msg)
	if e != nil {
		return nil, cli.NewCmdConfErr(e, nil)
	}

	return &Cmd{msg: msgReader, secretFile: opts.SecretFile}, nil
}

type Cmd struct {
	secretFile string
	msg        io.Reader
}

func (cmd *Cmd) Run() error {
	return digital_sign.Sign(cmd.msg, cmd.secretFile)
}
