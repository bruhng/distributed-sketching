package client

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/bruhng/distributed-sketching/proto"
	pb "github.com/bruhng/distributed-sketching/proto"
	"github.com/bruhng/distributed-sketching/sketches/kll"
)

func kllClient(k int, c pb.ServerClient) {

	sketch := kll.NewKLLSketch[int](k)
	i := 1

	for {
		// Read Data
		data := rand.Intn(100)
		sketch.Add(data)

		if i%100 == 0 {
			protoSketch := convertToOrdered2DArray(sketch.Sketch)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := c.MergeKll(ctx, protoSketch)
			if err != nil {
				panic(err)
			}
			sketch = kll.NewKLLSketch[int](k)

		}
		time.Sleep(200 * time.Millisecond)
		i++
	}
}

func convertToOrdered2DArray(data [][]int) *pb.Ordered2DArray {
	fmt.Println(data)
	orderedArray := &proto.Ordered2DArray{}

	for _, row := range data {
		protoRow := &proto.OrderedRow{} // Create a new row

		for _, val := range row {
			protoRow.Values = append(protoRow.Values, &proto.OrderedValue{
				Value: &proto.OrderedValue_IntVal{IntVal: int32(val)}, // Wrap value properly
			})
		}
		orderedArray.Rows = append(orderedArray.Rows, protoRow)
	}

	return orderedArray
}
