package client

import (
	"context"
	"math/rand"
	"time"

	pb "github.com/bruhng/distributed-sketching/proto"
	"github.com/bruhng/distributed-sketching/sketches/kll"
)

func kllClient(k int, c pb.SketcherClient) {

	sketch := kll.NewKLLSketch[int](k)
	i := 1

	for {
		// Read Data
		data := rand.Intn(100)
		sketch.Add(data)

		if i%10000 == 0 {
			protoSketch := convertToProtoKLL(sketch)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := c.MergeKll(ctx, protoSketch)
			if err != nil {
				panic(err)
			}
			sketch = kll.NewKLLSketch[int](k)

		}
		time.Sleep(10 * time.Microsecond)
		i++
	}
}

func convertToProtoKLL(sketch *kll.KLLSketch[int]) *pb.KLLSketch {
	orderedArray := &pb.KLLSketch{N: int32(sketch.N)}
	data := sketch.Sketch

	for _, row := range data {
		protoRow := &pb.OrderedRow{} // Create a new row

		for _, val := range row {
			protoRow.Values = append(protoRow.Values, &pb.OrderedValue{
				Value: &pb.OrderedValue_IntVal{IntVal: int32(val)}, // Wrap value properly
			})
		}
		orderedArray.Rows = append(orderedArray.Rows, protoRow)
	}

	return orderedArray
}
