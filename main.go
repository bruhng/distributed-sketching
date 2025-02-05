package main

import (
	"flag"

	"github.com/bruhng/distributed-sketching/client"
	"github.com/bruhng/distributed-sketching/server"
)

func main() {
	
	isClient := flag.Bool("client", false, "use flag if you want to create a client process instead of server")
	port := flag.String("port", "8080", "Choose what port to use")
	adress := flag.String("a", "127.0.0.1", "Choose what ip to connect to")
	flag.Parse()
	if *isClient {
		client.Init(*port, *adress)
	} else {
		server.Init(*port)
	}
}
