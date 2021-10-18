package recv

import (
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	"github.com/paulpaulych/crypto/internal/app/tcp"
	"github.com/paulpaulych/crypto/internal/infra/cli"
	"math/big"
	"net"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "recv"
}
func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	flagsSpec := cli.NewFlagSpec(conf.CmdName(), map[string]string{
		"host": "host to bind",
		"port": "port to bind",
		"P":    "large prime number",
		"Q":    "large prime number",
		"o":    "output type: file or console",
	})

	flags, err := flagsSpec.Parse(args)
	if err != nil {
		return nil, err
	}

	host := flags.Flags["host"].GetOr("localhost")
	port := flags.Flags["port"].GetOr("4444")
	addr := net.JoinHostPort(host, port)

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
	outputType := flags.Flags["o"].GetOr("console")
	output, e := cli.NewOutputFactory(outputType)
	if e != nil {
		return nil, cli.NewCmdConfError(e.Error(), nil)
	}

	if err != nil {
		return nil, err
	}
	return &Cmd{bindAddr: addr, output: output, p: P, q: Q}, nil
}

type Cmd struct {
	bindAddr string
	p, q     *big.Int
	output   cli.OutputFactory
}

func (cmd *Cmd) Run() error {
	return tcp.StartServer(cmd.bindAddr,
		msg_core.RecvMessage(
			protocols.RsaReader(cmd.p, cmd.q, cmd.output),
		),
	)
}
