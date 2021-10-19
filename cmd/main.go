package main

import (
	"fmt"
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/cmd/elgamal-msg"
	rsa_ds "github.com/paulpaulych/crypto/cmd/rsa-ds"
	"github.com/paulpaulych/crypto/cmd/rsa-msg"
	"github.com/paulpaulych/crypto/cmd/shamir-msg"
	"os"
	"strings"
)

func main() {

	subConfigs := []cli.CmdConf{
		&shamir_msg.Conf{},
		&elgamal_msg.Conf{},
		&rsa_msg.Conf{},
		&rsa_ds.Conf{},
	}

	cmd, confErr := cli.InitSubCmd(subConfigs, os.Args[1:])
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

func printConfError(e cli.CmdConfError) {
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
