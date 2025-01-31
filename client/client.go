package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/bruhng/distributed-sketching/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Init(port string, adr string) {
	conn, err := grpc.NewClient(adr+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		panic("Could not connect to server")
	}
	defer conn.Close()
	c := pb.NewServerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Merge(ctx, &pb.MergeRequest{Sketch: 12})
	if err != nil {
		fmt.Println("could not Merge")
	}
	fmt.Println("Response: ", r.GetStatus())

}
