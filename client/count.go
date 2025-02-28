package client

import (
	"context"
	"time"

	pb "github.com/bruhng/distributed-sketching/proto"
	"github.com/bruhng/distributed-sketching/shared"
	"github.com/bruhng/distributed-sketching/sketches/count"
	"github.com/bruhng/distributed-sketching/stream"
)

func countClient[T shared.Number](c pb.SketcherClient, dataStream stream.Stream[T]) {

	sketch := count.NewCountSketch[T](157, 100, 10)
	i := 1
	for {
		data := <-dataStream.Data
		sketch.Add(data)

		if i%10000 == 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			protoSketch := convertToProtoCount(sketch)
			_, err := c.MergeCount(ctx, protoSketch)
			if err != nil {
				panic(err)
			}
			sketch = count.NewCountSketch[T](157, 100, 10)

		}
		i++
	}
}

func convertToProtoCount[T shared.Number](sketch *count.CountSketch[T]) *pb.CountSketch {
	protoArray := &pb.CountSketch{}
	data := sketch.Sketch
	seeds := sketch.Seeds

	for _, row := range data {
		protoRow := &pb.IntRow{} // Create a new row

		for _, val := range row {
			protoRow.Val = append(protoRow.Val, int64(val))
		}
		protoArray.Rows = append(protoArray.Rows, protoRow)
	}

	for _, seed := range seeds {
		protoArray.Seeds = append(protoArray.Seeds, seed)
	}

	return protoArray
}
