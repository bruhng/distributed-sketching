package server

import (
	"context"

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
	kllState.Print()
	return &pb.MergeReply{Status: 0}, nil
}

func (s *server) QueryKll(_ context.Context, in *pb.OrderedValue) (*pb.QueryReturn, error) {
	val := in.GetIntVal()
	ret := kllState.Query(int(val))
	return &pb.QueryReturn{N: int32(kllState.N), Phi: int32(ret)}, nil
}

func (s *server) ReverseQueryKll(_ context.Context, in *pb.ReverseQuery) (*pb.OrderedValue, error) {
	phi := in.GetPhi()
	ret := kllState.QueryQuantile(float64(phi))
	return &pb.OrderedValue{Value: &pb.OrderedValue_IntVal{IntVal: int32(ret)}}, nil
}

func (s *server) PlotKll(_ context.Context, in *pb.OrderedValue) (*pb.PlotKllReply, error) {
	numBins := int(in.GetIntVal())
	xmin := kllState.QueryQuantile(0.0)
	xmax := kllState.QueryQuantile(1.0)
	step := float64(xmax-xmin) / float64(numBins)

	splits := make([]int, numBins+1)
	for i := 0; i <= numBins; i++ {
		splits[i] = xmin + int(step*float64(i))
	}

	pmf := make([]float32, numBins)
	for i := 0; i < numBins; i++ {
		pmf[i] = float32(kllState.Query(splits[i+1]) - kllState.Query(splits[i]))
	}
	return &pb.PlotKllReply{Step: float32(step), Pmf: pmf}, nil
}
