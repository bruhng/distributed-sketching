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
	sketchType := flag.String("sketch", "kll", "Choose what sketch to use")
	dataSetPath := flag.String("d", "", "Choose what data set path to use as data stream")
	dataSetName := flag.String("name", "", "Choose what part of the data set to use as data stream")
	dataSetType := flag.String("type", "int", "Choose what type the data set is")

	flag.Parse()
	if *isClient {
		if *dataSetType == "float" {
			client.Init[float64](*port, *adress, *sketchType, *dataSetPath, *dataSetName)
		} else if *dataSetType == "int" {
			client.Init[int](*port, *adress, *sketchType, *dataSetPath, *dataSetName)
		}
	} else if *isCustomer {
		customer.Init(*port, *adress)
	} else {
		server.Init(*port)
	}
}
