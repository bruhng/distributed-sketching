package client

import (
	"context"
	"fmt"
	"reflect"
	"time"

	pb "github.com/bruhng/distributed-sketching/proto"
	"github.com/bruhng/distributed-sketching/shared"
	"github.com/bruhng/distributed-sketching/sketches/kll"
	"github.com/bruhng/distributed-sketching/stream"
)

func kllClient[T shared.Number](k int, c pb.SketcherClient, dataStream stream.Stream[T]) {

	sketch := kll.NewKLLSketch[T](k)
	i := 1
	for {
		data := <-dataStream.Data

		sketch.Add(data)

		if i%100 == 0 {
			protoSketch := convertToProtoKLL(sketch)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := c.MergeKll(ctx, protoSketch)
			if err != nil {
				panic(err)
			}
			sketch = kll.NewKLLSketch[T](k)

		}
		i++

	}
}

func convertToProtoKLL[T shared.Number](sketch *kll.KLLSketch[T]) *pb.KLLSketch {
	t := fmt.Sprintf("%T", sketch.Sketch)[4:]
	orderedArray := &pb.KLLSketch{N: int32(sketch.N), Type: t}
	data := sketch.Sketch

	for _, row := range data {
		protoRow := &pb.OrderedRow{} // Create a new row

		for _, val := range row {
			if reflect.ValueOf(val).Kind() == reflect.Int {
				protoRow.Values = append(protoRow.Values, &pb.OrderedValue{
					Value: &pb.OrderedValue_IntVal{IntVal: int32(val)}, // Wrap value properly
				})

			} else {
				protoRow.Values = append(protoRow.Values, &pb.OrderedValue{
					Value: &pb.OrderedValue_FloatVal{FloatVal: float32(val)}, // Wrap value properly
				})

			}
		}
		orderedArray.Rows = append(orderedArray.Rows, protoRow)
	}

	return orderedArray
}
