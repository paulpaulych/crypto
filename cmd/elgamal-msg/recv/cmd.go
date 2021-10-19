package recv

import (
	"fmt"
	cli2 "github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	"github.com/paulpaulych/crypto/internal/app/tcp"
	dh "github.com/paulpaulych/crypto/internal/core/diffie-hellman"
	"math/big"
	"net"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "recv"
}
func (conf *Conf) NewCmd(args []string) (cli2.Cmd, cli2.CmdConfError) {
	flagsSpec := cli2.NewFlagSpec(conf.CmdName(), map[string]string{
		"host": "host to bind",
		"port": "port to bind",
		"P":    "large prime number",
		"G":    "generator of multiplicative group of integers modulo P",
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
	gStr := flags.Flags["G"].Get()
	if pStr == nil {
		return nil, cli2.NewCmdConfError("flag required: -P", nil)
	}
	if gStr == nil {
		return nil, cli2.NewCmdConfError("flag required: -G", nil)
	}
	P, success := new(big.Int).SetString(*pStr, 10)
	if !success {
		return nil, cli2.NewCmdConfError("cannot parse P", nil)
	}
	G, success := new(big.Int).SetString(*gStr, 10)
	if !success {
		return nil, cli2.NewCmdConfError("cannot parse G", nil)
	}
	commonPub, e := dh.NewCommonPublicKey(P, G)
	if e != nil {
		return nil, cli2.NewCmdConfError(
			fmt.Sprintf("Diffie-Hellman public key error: %v", e), nil,
		)
	}

	outputType := flags.Flags["o"].GetOr("console")
	output, e := cli2.NewOutputFactory(outputType)
	if e != nil {
		return nil, cli2.NewCmdConfError(e.Error(), nil)
	}

	if err != nil {
		return nil, err
	}
	return &Cmd{bindAddr: addr, output: output, commonPub: commonPub}, nil
}

type Cmd struct {
	bindAddr  string
	commonPub dh.CommonPublicKey
	output    cli2.OutputFactory
}

func (cmd *Cmd) Run() error {
	return tcp.StartServer(cmd.bindAddr,
		msg_core.RecvMessage(
			protocols.ElgamalReader(cmd.commonPub, cmd.output),
		),
	)
}
