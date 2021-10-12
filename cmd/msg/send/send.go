package send

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/nio"
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

	addr := flags.Args[0]
	msg := flags.Args[1]

	protocol := flags.Flags["protocol"].GetOr("shamir")
	prime := flags.Flags["prime"].Get()
	writer, err := writerForProtocol(protocol, prime)
	if err != nil {
		return nil, err
	}
	alreadyRead := new(uint)
	*alreadyRead = 0
	return &SendCmd{
		addr:  addr,
		alice: writer,
		msg:   &stringReader{s: msg, alreadyRead: alreadyRead},
	}, nil
}

type stringReader struct {
	s           string
	alreadyRead *uint
}

func (sr stringReader) TotalBytes() (uint, error) {
	return uint(len([]byte(sr.s))), nil
}

func (sr stringReader) Read(size uint) (*nio.BytePage, error) {
	bytes := []byte(sr.s)
	from := *sr.alreadyRead
	to := minUint(*sr.alreadyRead+size, uint(len(bytes)))
	*sr.alreadyRead = *(sr.alreadyRead) + size
	return &nio.BytePage{
		Bytes:   bytes[from:to],
		HasMore: to != uint(len(bytes)),
	}, nil
}

func minUint(x uint, y uint) uint {
	if x < y {
		return x
	} else {
		return y
	}
}

type SendCmd struct {
	addr  string
	alice msg_core.Alice
	msg   nio.ByteReader
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
