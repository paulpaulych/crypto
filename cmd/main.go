package main

import (
	"fmt"
	"github.com/paulpaulych/crypto/cmd/msg"
	cli2 "github.com/paulpaulych/crypto/internal/infra/cli"
	"os"
	"strings"
)

func main() {

	subConfigs := []cli2.CmdConf{
		&msg.MsgConf{},
	}

	cmd, confErr := cli2.InitSubCmd(subConfigs, os.Args[1:])
	if confErr != nil {
		printConfError(confErr)
		os.Exit(1)
	}

	runErr := cmd.Run()
	if runErr != nil {
		fmt.Println(runErr)
		os.Exit(1)
	}

	os.Exit(0)
}

func printConfError(e cli2.CmdConfError) {
	if e.Trace() != nil && len(e.Trace()) != 0 {
		path := strings.Join(e.Trace(), " ")
		fmt.Printf("%s: %s\n", path, e.Error())
	} else {
		fmt.Println(e.Error())
	}

	if e.Usage() == nil {
		return
	}
	println(*e.Usage())
}
