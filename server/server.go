package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	pb "github.com/bruhng/distributed-sketching/proto"
	"github.com/bruhng/distributed-sketching/sketches/kll"
)

type server struct {
	pb.UnimplementedSketcherServer
}

func newServer() *server {
	s := &server{}
	return s
}

func Init(port string) {

	kllState = kll.NewKLLSketch[int](200)

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(fmt.Sprint("listen error: ", err))
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSketcherServer(grpcServer, newServer())

	go grpcServer.Serve(l)

	for {

	}
}
