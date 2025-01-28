package server

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/bruhng/distributed-sketching/types"
)

type Server struct {
	Sketch types.Sketch
}

func (s *Server) Merge(args *types.Args, reply *types.Reply) error {
	fmt.Printf("wow merge is called")
	*reply = 69
	return nil
}

func Init(port string) {
	server := new(Server)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(fmt.Sprint("listen error: ", err))
	}
	go http.Serve(l, nil)

	for {

	}
}
