package recv

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	"github.com/paulpaulych/crypto/internal/app/tcp"
	"net"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "recv"
}
func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	var opts struct {
		Host   string `short:"h" long:"host" description:"host to bind" default:"localhost"`
		Port   string `short:"p" long:"port" description:"port to bind" default:"12345"`
		Output string `short:"o" long:"output" choice:"file" choice:"console" description:"output type: console or file" default:"console"`
	}

	_, err := cli.ParseFlagsOfCmd(conf.CmdName(), &opts, args)
	if err != nil {
		return nil, err
	}

	output, e := cli.NewOutputFactory(opts.Output)
	if e != nil {
		return nil, cli.NewCmdConfErr(e, nil)
	}

	return &Cmd{
		bindAddr: net.JoinHostPort(opts.Host, opts.Port),
		output:   output,
	}, nil
}

type Cmd struct {
	bindAddr string
	output   cli.OutputFactory
}

func (cmd *Cmd) Run() error {
	return tcp.StartServer(cmd.bindAddr, msg_core.RecvMessage(protocols.ShamirReader(cmd.output)))
}
