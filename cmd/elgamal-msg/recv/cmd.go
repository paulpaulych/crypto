package recv

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	"github.com/paulpaulych/crypto/internal/app/tcp"
	"math/big"
	"net"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "recv"
}
func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	var opts struct {
		Host   string         `short:"h" long:"host" description:"host to bind" default:"localhost"`
		Port   string         `short:"p" long:"port" description:"port to bind" default:"12345"`
		P      *cli.BigIntOpt `short:"P" description:"large prime number" required:"true"`
		G      *cli.BigIntOpt `short:"G" description:"generator of multiplicative group of integers modulo P" required:"true"`
		Output string         `short:"o" long:"output" choice:"file" choice:"console" description:"output type: console or file" default:"console"`
	}
	if _, err := cli.ParseFlagsOfCmd(conf.CmdName(), &opts, args); err != nil {
		return nil, err
	}

	output, e := cli.NewOutputFactory(opts.Output)
	if e != nil {
		return nil, cli.NewCmdConfErr(e, nil)
	}

	return &Cmd{
		bindAddr: net.JoinHostPort(opts.Host, opts.Port),
		output:   output,
		p:        opts.P.Value,
		g:        opts.G.Value,
	}, nil
}

type Cmd struct {
	bindAddr string
	p, g     *big.Int
	output   cli.OutputFactory
}

func (cmd *Cmd) Run() error {
	reader, err := protocols.ElgamalReader(cmd.p, cmd.g, cmd.output)
	if err != nil {
		return err
	}
	return tcp.StartServer(cmd.bindAddr,
		msg_core.RecvMessage(reader),
	)
}
