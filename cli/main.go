package main

import (
	"fmt"
	. "github.com/paulpaulych/crypto/shamir/cli/cmd/common"
	"github.com/paulpaulych/crypto/shamir/cli/cmd/msg"
	"os"
	"strings"
)

func main() {

	subConfigs := []CmdConf{
		&msg.MsgConf{},
	}

	cmd, confErr := InitSubCmd(subConfigs, os.Args[1:])
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

func printConfError(e CmdConfError) {
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
