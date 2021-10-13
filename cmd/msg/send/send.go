package send

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	"github.com/paulpaulych/crypto/internal/infra/cli"
	"io"
	"math/big"
	"os"
	"strings"
)

type SendConf struct{}

func (conf *SendConf) CmdName() string {
	return "send"
}

func (conf *SendConf) InitCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	flagsSpec := cli.NewFlagSpec(conf.CmdName(), map[string]string{
		"protocol": "protocol",
		"prime":    "prime integer",
		"i":        "message input type: file or arg",
	})

	flags, err := flagsSpec.Parse(args)
	if err != nil {
		return nil, err
	}

	if len(flags.Args) < 2 {
		return nil, cli.NewCmdConfError("args required: [host:port] [message]", nil)
	}

	addr := flags.Args[0]
	msg := flags.Args[1]

	protocol := flags.Flags["protocol"].GetOr("shamir")
	prime := flags.Flags["prime"].Get()
	input := flags.Flags["i"].GetOr("console")
	msgReader, e := chooseMsgReader(input, msg)
	if e != nil {
		return nil, cli.NewCmdConfError(e.Error(), nil)
	}
	writer, e := writerForProtocol(protocol, prime)
	if err != nil {
		return nil, err
	}
	return &SendCmd{addr: addr, alice: writer, msg: msgReader}, nil
}

func chooseMsgReader(input string, msg string) (io.Reader, error) {
	switch input {
	case "console":
		return strings.NewReader(msg), nil
	case "file":
		f, err := os.Open(msg)
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		return nil, errors.New("unknown input type")
	}
}

type SendCmd struct {
	addr  string
	alice msg_core.Alice
	msg   io.Reader
}

func (cmd *SendCmd) Run() error {
	return msg_core.SendMsg(cmd.addr, cmd.msg, cmd.alice)
}

func writerForProtocol(name string, primeStr *string) (msg_core.Alice, cli.CmdConfError) {
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
