package elgamal_msg

import (
	"github.com/paulpaulych/crypto/cmd/elgamal-msg/recv"
	"github.com/paulpaulych/crypto/cmd/elgamal-msg/send"
	"github.com/paulpaulych/crypto/internal/infra/cli"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "elgamal-msg"
}

func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	subConfigs := []cli.CmdConf{
		&send.Conf{},
		&recv.Conf{},
	}
	return cli.InitSubCmd(subConfigs, args)
}
