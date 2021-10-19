package elgamal_msg

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/cmd/elgamal-msg/recv"
	"github.com/paulpaulych/crypto/cmd/elgamal-msg/send"
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
