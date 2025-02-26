package client

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/bruhng/distributed-sketching/proto"
	"github.com/bruhng/distributed-sketching/sketches/count"
	"github.com/bruhng/distributed-sketching/stream"
)

func countClient(c pb.SketcherClient, dataSetPath string) {

	sketch := count.NewCountSketch[int](157, 100, 10)
	i := 1
	dataStream := stream.NewStreamFromPath(dataSetPath)
	for {
		strData := <-dataStream.Data
		if strData == "" {
			continue
		}
		for _, char := range strData {
			fmt.Println(int(char))
		}
		data, err := strconv.Atoi(strData)

		if err != nil {
			fmt.Println(err)
			panic("Could not convert file to int")
		}
		sketch.Add(data)

		if i%10000 == 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			protoSketch := convertToProtoCount(sketch)
			_, err := c.MergeCount(ctx, protoSketch)
			if err != nil {
				panic(err)
			}
			sketch = count.NewCountSketch[int](157, 100, 10)

		}
		i++
	}
}

func convertToProtoCount(sketch *count.CountSketch[int]) *pb.CountSketch {
	protoArray := &pb.CountSketch{}
	data := sketch.Sketch
	seeds := sketch.Seeds

	for _, row := range data {
		protoRow := &pb.IntRow{} // Create a new row

		for _, val := range row {
			protoRow.Val = append(protoRow.Val, int32(val))
		}
		protoArray.Rows = append(protoArray.Rows, protoRow)
	}

	for _, seed := range seeds {
		protoArray.Seeds = append(protoArray.Seeds, seed)
	}

	return protoArray
}
