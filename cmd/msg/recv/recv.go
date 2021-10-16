package recv

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/nio"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	"github.com/paulpaulych/crypto/internal/app/tcp"
	"github.com/paulpaulych/crypto/internal/infra/cli"
	"log"
	"net"
	"time"
)

type RecvConf struct{}

func (conf *RecvConf) CmdName() string {
	return "recv"
}

func (conf *RecvConf) InitCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	flagsSpec := cli.NewFlagSpec(conf.CmdName(), map[string]string{
		"host": "host to bind",
		"port": "port to bind",
		"o":    "output type: file or console",
	})

	flags, err := flagsSpec.Parse(args)
	if err != nil {
		return nil, err
	}

	host := flags.Flags["host"].GetOr("localhost")
	port := flags.Flags["port"].GetOr("4444")
	addr := net.JoinHostPort(host, port)

	outputType := flags.Flags["o"].GetOr("console")
	output, e := outputFactory(outputType)
	if e != nil {
		return nil, cli.NewCmdConfError(e.Error(), nil)
	}

	if err != nil {
		return nil, err
	}
	return &RecvCmd{bindAddr: addr, output: output}, nil
}

type RecvCmd struct {
	bindAddr string
	output   OutputFactory
}

func (cmd *RecvCmd) Run() error {
	onErr := func(e error) {
		log.Printf("error reading message: %s", e)
	}
	chooseBob := func(code msg_core.ProtocolCode) (msg_core.Bob, error) {
		return protocols.ChooseBob(code, cmd.output, onErr)
	}
	return tcp.StartServer(cmd.bindAddr, msg_core.RecvMessage(chooseBob))
}

type OutputFactory = func(from net.Addr) nio.ClosableWriter

func outputFactory(output string) (OutputFactory, error) {
	switch output {
	case "console":
		return func(from net.Addr) nio.ClosableWriter {
			prefix := fmt.Sprintln("RECEIVED MESSAGE FROM", from)
			return nio.NewConsoleWriter([]byte(prefix))
		}, nil
	case "file":
		return func(from net.Addr) nio.ClosableWriter {
			fName := fmt.Sprintf("msg_from_%s_at_%v.txt", from.String(), time.Now().UnixMilli())
			return nio.NewFileWriter(fName, func() {
				log.Printf("received message saved to file %s", fName)
			})
		}, nil
	default:
		return nil, errors.New("invalid output type")
	}
}
