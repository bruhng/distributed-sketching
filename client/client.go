package client

import (
	"fmt"

	pb "github.com/bruhng/distributed-sketching/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Init(port string, adr string, sketchType string) {
	conn, err := grpc.NewClient(adr+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		panic("Could not connect to server")
	}
	defer conn.Close()
	c := pb.NewServerClient(conn)

	switch sketchType {
	case "kll":
		kllClient(100, c)
	default:
		panic("No sketch provided or invalid sketch")
	}
}
