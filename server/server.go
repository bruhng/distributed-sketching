package server

import (
	"fmt"
	"net"

	pb "github.com/bruhng/distributed-sketching/proto"
	"github.com/bruhng/distributed-sketching/sketches/count"
	"github.com/bruhng/distributed-sketching/sketches/kll"
	"google.golang.org/grpc"
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
	countState = count.NewCountSketch[int](157, 100, 10)

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
