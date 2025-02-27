package client

import (
	"fmt"

	pb "github.com/bruhng/distributed-sketching/proto"
	"github.com/bruhng/distributed-sketching/shared"
	"github.com/bruhng/distributed-sketching/stream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Init[T shared.Number](port string, adr string, sketchType string, dataSetPath string, headerName string) {
	conn, err := grpc.NewClient(adr+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		panic("Could not connect to server")
	}
	defer conn.Close()
	c := pb.NewSketcherClient(conn)
	dataStream := *stream.NewStreamFromCsv[T](dataSetPath, headerName)

	switch sketchType {
	case "kll":
		kllClient(100, c, dataStream)
	case "count":
		countClient(c, dataStream)
	default:
		panic("No sketch provided or invalid sketch")
	}
}
