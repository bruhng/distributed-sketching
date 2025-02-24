package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	pb "github.com/bruhng/distributed-sketching/proto"
)

type server struct {
	pb.UnimplementedServerServer
}

func newServer() *server {
	s := &server{}
	return s
}

func (s *server) MergeKll(_ context.Context, in *pb.Ordered2DArray) (*pb.MergeReply, error) {
	fmt.Println("arr", in.Rows)

	return &pb.MergeReply{Status: 0}, nil
}

func Init(port string) {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(fmt.Sprint("listen error: ", err))
	}
	grpcServer := grpc.NewServer()
	pb.RegisterServerServer(grpcServer, newServer())

	go grpcServer.Serve(l)

	for {

	}
}
