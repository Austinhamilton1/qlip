package main

import (
	"flag"
	"fmt"

	"github.com/Austinhamilton1/qlip/internal/qlipnet"
)

func main() {
	portPtr := flag.Int("port", 8080, "Port to use for networking.")
	remoteIPPtr := flag.String("addr", "localhost", "Remote IP to connect to.")

	flag.Parse()

	typ := flag.Args()[0]

	switch typ {
	case "server":
		server := qlipnet.NewServer()
		server.Serve(*portPtr)
	case "client":
		client := qlipnet.NewClient()
		err := client.Connect(*remoteIPPtr, *portPtr)
		if err != nil {
			fmt.Printf("could not connect to server: %v\n", err)
			return
		}
		client.Loop()
	default:
		fmt.Println("invalid type: ", typ)
	}
}
