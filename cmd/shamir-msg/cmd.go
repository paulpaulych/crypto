package shamir_msg

import (
	"github.com/paulpaulych/crypto/cmd/shamir-msg/recv"
	"github.com/paulpaulych/crypto/cmd/shamir-msg/send"
	"github.com/paulpaulych/crypto/internal/infra/cli"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "shamir-msg"
}

func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	subConfigs := []cli.CmdConf{
		&send.Conf{},
		&recv.Conf{},
	}
	return cli.InitSubCmd(subConfigs, args)
}
