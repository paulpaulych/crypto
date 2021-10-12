package send

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	"github.com/paulpaulych/crypto/internal/infra/cli"
	"math/big"
)

type SendConf struct{}

func (conf *SendConf) CmdName() string {
	return "send"
}

func (conf *SendConf) InitCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	flagsSpec := cli.NewFlagSpec(conf.CmdName(), map[string]string{
		"protocol": "protocol",
		"prime":    "prime integer",
	})

	flags, err := flagsSpec.Parse(args)
	if err != nil {
		return nil, err
	}

	if len(flags.Args) < 2 {
		return nil, cli.NewCmdConfError("args required: [host:port] [message]", nil)
	}

	addr, msgStr := flags.Args[0], flags.Args[1]
	msg, success := new(big.Int).SetString(msgStr, 10)
	if !success {
		return nil, cli.NewCmdConfError("message must be integer", nil)
	}

	protocol := flags.Flags["protocol"].GetOr("shamir")
	prime := flags.Flags["prime"].Get()
	writer, err := writerForProtocol(protocol, prime)
	if err != nil {
		return nil, err
	}
	return &SendCmd{addr: addr, writer: writer, msg: msg}, nil
}

type SendCmd struct {
	addr   string
	writer msg_core.MsgWriter
	msg    msg_core.Msg
}

func (cmd *SendCmd) Run() error {
	return msg_core.SendMsg(cmd.addr, cmd.msg, cmd.writer)
}

func writerForProtocol(name string, primeStr *string) (msg_core.MsgWriter, cli.CmdConfError) {
	switch name {
	case "shamir":
		if primeStr == nil || len(*primeStr) == 0 {
			return nil, cli.NewCmdConfError("shamir protocol requires -prime flag", nil)
		}
		prime, success := new(big.Int).SetString(*primeStr, 10)
		if !success {
			return nil, cli.NewCmdConfError("cannot parse prime", nil)
		}
		return protocols.ShamirWriter(prime), nil
	default:
		msg := fmt.Sprintf("unknown protocol '%s'", name)
		return nil, cli.NewCmdConfError(msg, nil)
	}
}
