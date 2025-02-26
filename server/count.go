package server

import (
	"context"
	"fmt"

	pb "github.com/bruhng/distributed-sketching/proto"
	"github.com/bruhng/distributed-sketching/sketches/count"
)

var countState *count.CountSketch[int]

func covertProtoCountToCount(protoData *pb.CountSketch) *count.CountSketch[int] {
	var data [][]int
	var seeds []uint32

	for _, protoRow := range protoData.Rows {
		var row []int

		for _, protoValue := range protoRow.Val {
			row = append(row, int(protoValue))
		}

		data = append(data, row)
	}
	for _, protoSeed := range protoData.Seeds {
		seeds = append(seeds, protoSeed)
	}

	return count.NewCountFromData[int](data, seeds)
}

func (s *server) MergeCount(_ context.Context, in *pb.CountSketch) (*pb.MergeReply, error) {
	sketch := covertProtoCountToCount(in)
	countState.Merge(*sketch)
	fmt.Println(countState)
	return &pb.MergeReply{Status: 0}, nil
}

func (s *server) QueryCount(_ context.Context, in *pb.AnyValue) (*pb.CountQueryReply, error) {
	val := in.GetIntVal()
	ret := countState.Query(int(val))
	return &pb.CountQueryReply{Res: int32(ret)}, nil
}
