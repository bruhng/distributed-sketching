package main

import (
	// "flag"
	"fmt"
	"math/rand"

	// "github.com/bruhng/distributed-sketching/client"
	// "github.com/bruhng/distributed-sketching/server"
	"github.com/bruhng/distributed-sketching/sketches/kll"
)

func main() {

	sketch := kll.NewKLLSketch[int, int](200)

	for i := 0; i <= 10000; i++ {
		sketch.Add(rand.Intn(100))
	}
	fmt.Println(sketch.Query(0))
	fmt.Println(sketch.Query(50))
	fmt.Println(sketch.Query(99))
	fmt.Println(sketch.QueryQuantile(5000))

	// isClient := flag.Bool("client", false, "use flag if you want to create a client process instead of server")
	// port := flag.String("port", "8080", "Choose what port to use")
	// adress := flag.String("a", "127.0.0.1", "Choose what ip to connect to")
	// flag.Parse()
	// if *isClient {
	// 	client.Init(*port, *adress)
	// } else {
	// 	server.Init(*port)
	// }
}
