package key_gen

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	digital_sign "github.com/paulpaulych/crypto/internal/app/digital-sign"
	"math/big"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "key-gen"
}
func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	flagsSpec := cli.NewFlagSpec(conf.CmdName(), map[string]string{
		"P": "large prime number",
		"Q": "large prime number",
	})

	flags, err := flagsSpec.Parse(args)
	if err != nil {
		return nil, err
	}

	pStr := flags.Flags["P"].Get()
	qStr := flags.Flags["Q"].Get()
	if pStr == nil {
		return nil, cli.NewCmdConfError("flag required: -P", nil)
	}
	if qStr == nil {
		return nil, cli.NewCmdConfError("flag required: -Q", nil)
	}
	P, success := new(big.Int).SetString(*pStr, 10)
	if !success {
		return nil, cli.NewCmdConfError("cannot parse P", nil)
	}
	Q, success := new(big.Int).SetString(*qStr, 10)
	if !success {
		return nil, cli.NewCmdConfError("cannot parse Q", nil)
	}

	if err != nil {
		return nil, err
	}
	return &Cmd{p: P, q: Q}, nil
}

type Cmd struct {
	p, q   *big.Int
	output cli.OutputFactory
}

func (cmd *Cmd) Run() error {
	return digital_sign.GenerateKeys(cmd.p, cmd.q)
}
