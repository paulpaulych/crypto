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
	if e := confErr; e != nil {
		if e.Trace() != nil && len(e.Trace()) != 0 {
			path := strings.Join(e.Trace(), " ")
			fmt.Printf("%s:\n", path)
		}
		if writeHelp := e.HelpWriter(); writeHelp != nil {
			writeHelp(os.Stdout)
		} else if err := e.Error(); err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}
		os.Exit(1)
	}

	runErr := cmd.Run()
	if runErr != nil {
		fmt.Println(runErr)
		os.Exit(1)
	}

	os.Exit(0)
}
