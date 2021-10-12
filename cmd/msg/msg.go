package msg

import (
	"github.com/paulpaulych/crypto/cmd/msg/recv"
	"github.com/paulpaulych/crypto/cmd/msg/send"
	cli2 "github.com/paulpaulych/crypto/internal/infra/cli"
)

type MsgConf struct{}

func (conf *MsgConf) CmdName() string {
	return "msg"
}

func (conf *MsgConf) InitCmd(args []string) (cli2.Cmd, cli2.CmdConfError) {
	subConfigs := []cli2.CmdConf{
		&send.SendConf{},
		&recv.RecvConf{},
	}
	return cli2.InitSubCmd(subConfigs, args)
}
