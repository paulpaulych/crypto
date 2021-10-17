package send

import (
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	"github.com/paulpaulych/crypto/internal/infra/cli"
	"io"
	"math/big"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "send"
}

func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	flagsSpec := cli.NewFlagSpec(conf.CmdName(), map[string]string{
		"P": "large prime integer",
		"i": "message input type: file or arg",
	})

	flags, err := flagsSpec.Parse(args)
	if err != nil {
		return nil, err
	}

	if len(flags.Args) < 2 {
		return nil, cli.NewCmdConfError("args required: [host:port] [message]", nil)
	}

	addr := flags.Args[0]

	primeStr := flags.Flags["P"].Get()
	input := flags.Flags["i"].GetOr("console")

	msgReader, e := cli.NewInputReader(input, flags.Args[1:])
	if e != nil {
		return nil, cli.NewCmdConfError(e.Error(), nil)
	}

	if primeStr == nil || len(*primeStr) == 0 {
		return nil, cli.NewCmdConfError("shamir protocol requires -prime flag", nil)
	}
	prime, success := new(big.Int).SetString(*primeStr, 10)
	if !success {
		return nil, cli.NewCmdConfError("cannot parse P", nil)
	}
	if err != nil {
		return nil, err
	}
	return &Cmd{addr: addr, alice: protocols.ShamirWriter(prime), msg: msgReader}, nil
}

type Cmd struct {
	addr  string
	alice msg_core.ConnWriter
	msg   io.Reader
}

func (cmd *Cmd) Run() error {
	return msg_core.SendMsg(cmd.addr, cmd.msg, cmd.alice)
}
