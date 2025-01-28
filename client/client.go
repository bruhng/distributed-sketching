package client

import (
	"fmt"
	"net/rpc"

	"github.com/bruhng/distributed-sketching/types"
)

func Init(port string, adr string) {
	client, err := rpc.DialHTTP("tcp", adr+":"+port)
	if err != nil {
		panic(fmt.Sprint("Dial error:", client))
	}

	args := types.Args{Sketch: 1}
	var reply types.Reply
	err = client.Call("Server.Merge", args, &reply)
	if err != nil {
		fmt.Println("ohh no", err)
	}
	fmt.Println(reply)
	fmt.Println("done")

}
