package main

import (
	"os"

	"github.com/bruhng/distributed-sketching/client"
	"github.com/bruhng/distributed-sketching/server"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("Pleas enter client or server")
	}
	runType := args[0]
	if runType == "client" {
		client.Init("8080")
	} else if runType == "server" {
		server.Init("8080")
	} else {
		panic("Not a valid run type, pleas enter client or server")
	}
}
