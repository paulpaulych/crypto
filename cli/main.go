package main

import (
	"flag"
	"fmt"
	"github.com/paulpaulych/crypto/shamir/cli/cmd"
	args2 "github.com/paulpaulych/crypto/shamir/cli/cmd/flags"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("send or recv subcommand is required")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "send":
		flagSet := flag.NewFlagSet("send", flag.ExitOnError)
		var flags, err = args2.ParseSendCmdFlags(flagSet, os.Args[2:])
		if err != nil {
			flag.PrintDefaults()
			os.Exit(1)
		}
		err = command.SendShamirEncrypted(flags)
		if err != nil {
			log.Printf("can't send: %v", err)
			os.Exit(1)
		}
	case "recv":
		flagSet := flag.NewFlagSet("recv", flag.ExitOnError)
		flags, err := args2.ParseRecvCmdFlags(flagSet, os.Args[2:])
		if err != nil {
			flag.PrintDefaults()
			os.Exit(1)
		}
		err = command.RecvMessages(flags)
		if err != nil {
			log.Printf("can't start server: %v", err)
			return
		}
	default:
		fmt.Println("only send and recv commands supported")
		os.Exit(1)
	}
}
