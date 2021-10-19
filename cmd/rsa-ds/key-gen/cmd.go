package key_gen

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app/digital-sign"
	"math/big"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "key-gen"
}

func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	var opts struct {
		P *cli.BigIntOpt `short:"p" description:"large prime number" required:"true"`
		Q *cli.BigIntOpt `short:"q" description:"large prime number" required:"true"`
	}
	_, err := cli.ParseFlagsOfCmd(conf.CmdName(), &opts, args)
	if err != nil {
		return nil, err
	}
	return &Cmd{p: opts.P.Value, q: opts.Q.Value}, nil
}

type Cmd struct {
	p, q *big.Int
}

func (cmd *Cmd) Run() error {
	return digital_sign.GenerateKeys(cmd.p, cmd.q)
}
