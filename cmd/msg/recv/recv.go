package recv

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/messaging"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	"github.com/paulpaulych/crypto/internal/infra/cli"
	"net"
)

type RecvConf struct{}

func (conf *RecvConf) CmdName() string {
	return "recv"
}

func (conf *RecvConf) InitCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	flagsSpec := cli.NewFlagSpec(conf.CmdName(), map[string]string{
		"protocol": "protocol",
		"host":     "host to bind",
		"port":     "port to bind",
	})

	flags, err := flagsSpec.Parse(args)
	if err != nil {
		return nil, err
	}

	host := flags.Flags["host"].GetOr("localhost")
	port := flags.Flags["port"].GetOr("4444")
	addr := net.JoinHostPort(host, port)

	protocol := flags.Flags["protocol"].GetOr("shamir")
	reader, err := readerForProtocol(protocol)
	if err != nil {
		return nil, err
	}
	return &RecvCmd{bindAddr: addr, reader: reader}, nil
}

type RecvCmd struct {
	bindAddr string
	reader   messaging.ReadMsg
}

func (cmd *RecvCmd) Run() error {
	return messaging.ListenForMsg(cmd.bindAddr, cmd.reader)
}

func readerForProtocol(name string) (messaging.ReadMsg, cli.CmdConfError) {
	switch name {
	case "shamir":
		return protocols.ShamirReader(), nil
	default:
		msg := fmt.Sprintf("unknown protocol '%s'", name)
		return nil, cli.NewCmdConfError(msg, nil)
	}
}
