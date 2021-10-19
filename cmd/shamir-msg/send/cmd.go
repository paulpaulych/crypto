package send

import (
	cli2 "github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	"io"
	"math/big"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "send"
}

func (conf *Conf) NewCmd(args []string) (cli2.Cmd, cli2.CmdConfError) {
	flagsSpec := cli2.NewFlagSpec(conf.CmdName(), map[string]string{
		"P": "large prime integer",
		"i": "message input type: file or arg",
	})

	flags, err := flagsSpec.Parse(args)
	if err != nil {
		return nil, err
	}

	if len(flags.Args) < 2 {
		return nil, cli2.NewCmdConfError("args required: [host:port] [message]", nil)
	}

	addr := flags.Args[0]

	primeStr := flags.Flags["P"].Get()
	input := flags.Flags["i"].GetOr("console")

	msgReader, e := cli2.NewInputReader(input, flags.Args[1:])
	if e != nil {
		return nil, cli2.NewCmdConfError(e.Error(), nil)
	}

	if primeStr == nil || len(*primeStr) == 0 {
		return nil, cli2.NewCmdConfError("shamir protocol requires -prime flag", nil)
	}
	prime, success := new(big.Int).SetString(*primeStr, 10)
	if !success {
		return nil, cli2.NewCmdConfError("cannot parse P", nil)
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
