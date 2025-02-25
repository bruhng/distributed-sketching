package server

import (
	"context"
	"fmt"

	pb "github.com/bruhng/distributed-sketching/proto"
	"github.com/bruhng/distributed-sketching/sketches/kll"
)

var kllState *kll.KLLSketch[int]

func covertProtoKLLToKLL(protoData *pb.KLLSketch) *kll.KLLSketch[int] {
	var data [][]int

	for _, protoRow := range protoData.Rows {
		var row []int

		for _, protoValue := range protoRow.Values {
			if intVal, ok := protoValue.Value.(*pb.OrderedValue_IntVal); ok {
				row = append(row, int(intVal.IntVal)) // Convert back to int
			}
		}

		data = append(data, row)
	}

	return kll.NewKLLFromData[int](data, int(protoData.GetN()), 200)
}

func (s *server) MergeKll(_ context.Context, in *pb.KLLSketch) (*pb.MergeReply, error) {
	sketch := covertProtoKLLToKLL(in)
	kllState.Merge(*sketch)
	fmt.Println("Recived sketch, state is now: ", kllState)
	return &pb.MergeReply{Status: 0}, nil
}
