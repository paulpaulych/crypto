package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"shamir-cli/args"
)

const maxMsgBytes = 2

func main() {
	if len(os.Args) < 2 {
		fmt.Println("send or recv subcommand is required")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "send":
		flagSet := flag.NewFlagSet("send", flag.ExitOnError)
		var flags, err = args.ParseSendCmdFlags(flagSet, os.Args[2:])
		if err != nil {
			flag.PrintDefaults()
			os.Exit(1)
		}
		bytes := toBytes(flags.Msg)
		err = tcpSend(flags.Addr, bytes)
		if err != nil {
			log.Printf("can't send: %v", err)
			os.Exit(1)
		}
	case "recv":
		flagSet := flag.NewFlagSet("recv", flag.ExitOnError)
		flags, err := args.ParseRecvCmdFlags(flagSet, os.Args[2:])
		if err != nil {
			flag.PrintDefaults()
			os.Exit(1)
		}
		err = startServer(flags.BindAddr, maxMsgBytes, printReceivedBytes)
		if err != nil {
			log.Printf("can't start server: %v", err)
			return
		}
	default:
		fmt.Println("only send and recv commands supported")
		os.Exit(1)
	}
}

func toBytes(msg string) []byte {
	bigint := new(big.Int)
	bigint.SetString(msg, 10)
	bytes := bigint.Bytes()
	log.Printf("size: %v, bytes: %v", len(bytes), bytes)
	log.Println()
	return bytes
}

func printReceivedBytes(bytes []byte) {
	res := new(big.Int)
	res.SetBytes(bytes)
	fmt.Printf("received string: %s", res)
	fmt.Println()
}
