package main

import (
	"flag"
	// "fmt"
	// "math/rand"

	"github.com/bruhng/distributed-sketching/client"
	"github.com/bruhng/distributed-sketching/customer"
	"github.com/bruhng/distributed-sketching/server"
	// "github.com/bruhng/distributed-sketching/sketches/kll"
)

func main() {

	isClient := flag.Bool("client", false, "use flag if you want to create a client process instead of server")
	isCustomer := flag.Bool("customer", false, "use flag if you want to create a customer process instead of a server")
	port := flag.String("port", "8080", "Choose what port to use")
	adress := flag.String("a", "127.0.0.1", "Choose what ip to connect to")
	flag.Parse()
	if *isClient {
		client.Init(*port, *adress, "kll")
	} else if *isCustomer {
		customer.Init(*port, *adress)
	} else {
		server.Init(*port)
	}
}
