package server

import (
	"cmp"
	"context"
	"fmt"
	"sync"

	pb "github.com/bruhng/distributed-sketching/proto"
	"github.com/bruhng/distributed-sketching/sketches/kll"
)

var (
	kllStateOnce sync.Once
	kllStateMap  sync.Map
)

func getOrCreateKllState[T cmp.Ordered]() *kll.KLLSketch[T] {
	key := fmt.Sprintf("%T", *new(T))
	var sketch *kll.KLLSketch[T]
	kllStateOnce.Do(func() {
		kllStateMap.Store(key, kll.NewKLLSketch[T](200))
	})
	if val, ok := kllStateMap.Load(key); ok {
		sketch = val.(*kll.KLLSketch[T])
	}
	return sketch
}

func covertProtoKLLToKLL[T cmp.Ordered](protoData *pb.KLLSketch) *kll.KLLSketch[T] {
	var data [][]T

	for _, protoRow := range protoData.Rows {
		var row []T

		for _, protoValue := range protoRow.Values {
			if intVal, ok := protoValue.Value.(*pb.OrderedValue_IntVal); ok {
				row = append(row, T(intVal.IntVal))
			}
			if floatVal, ok := protoValue.Value.(*pb.OrderedValue_FloatVal); ok {
				row = append(row, T(any(floatVal.FloatVal).(T)))
			}
		}

		data = append(data, row)
	}

	return kll.NewKLLFromData[T](data, int(protoData.GetN()), 200)
}

func (s *server) MergeKll(_ context.Context, in *pb.KLLSketch) (*pb.MergeReply, error) {

	if in.Type == "int" {
		kllState := getOrCreateKllState[int]()
		sketch := covertProtoKLLToKLL[int](in)
		kllState.Merge(*sketch)
		kllState.Print()
	} else if in.Type == "float64" {
		kllState := getOrCreateKllState[float32]()
		sketch := covertProtoKLLToKLL[float32](in)
		kllState.Merge(*sketch)
		kllState.Print()
	} else {
		return nil, fmt.Errorf("Type submitted is not supported")
	}

	return &pb.MergeReply{Status: 0}, nil
}

func (s *server) QueryKll(_ context.Context, in *pb.OrderedValue) (*pb.QueryReturn, error) {

	if in.Type == "int" {
		kllState := getOrCreateKllState[int]()
		val := in.GetIntVal()
		ret := kllState.Query(int(val))
		return &pb.QueryReturn{N: int64(kllState.N), Phi: int64(ret)}, nil
	} else if in.Type == "float64" {
		kllState := getOrCreateKllState[float32]()
		val := in.GetIntVal()
		ret := kllState.Query(float32(val))
		return &pb.QueryReturn{N: int64(kllState.N), Phi: int64(ret)}, nil
	} else {
		return nil, fmt.Errorf("Type submitted is not supported")
	}
}

func (s *server) ReverseQueryKll(_ context.Context, in *pb.ReverseQuery) (*pb.OrderedValue, error) {
	phi := in.Phi
	if in.Type == "int" {
		kllState := getOrCreateKllState[int]()
		ret := kllState.QueryQuantile(float64(phi))
		return &pb.OrderedValue{Value: &pb.OrderedValue_IntVal{IntVal: int64(ret)}}, nil
	} else if in.Type == "float64" {
		kllState := getOrCreateKllState[float32]()
		ret := kllState.QueryQuantile(float64(phi))
		return &pb.OrderedValue{Value: &pb.OrderedValue_FloatVal{FloatVal: float32(ret)}}, nil
	} else {
		return nil, fmt.Errorf("Type submitted is not supported")
	}
}

func (s *server) PlotKll(_ context.Context, in *pb.PlotRequest) (*pb.PlotKllReply, error) {
	if in.Type == "int" {
		kllState := getOrCreateKllState[int]()
		numBins := int(in.GetNumBins())
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
	} else if in.Type == "float64" {
		kllState := getOrCreateKllState[float32]()
		numBins := int(in.GetNumBins())
		xmin := kllState.QueryQuantile(0.0)
		xmax := kllState.QueryQuantile(1.0)
		step := float32(xmax-xmin) / float32(numBins)

		splits := make([]float32, numBins+1)
		for i := 0; i <= numBins; i++ {
			splits[i] = xmin + step*float32(i)
		}

		pmf := make([]float32, numBins)
		for i := 0; i < numBins; i++ {
			pmf[i] = float32(kllState.Query(splits[i+1]) - kllState.Query(splits[i]))
		}
		return &pb.PlotKllReply{Step: float32(step), Pmf: pmf}, nil

	} else {
		return nil, fmt.Errorf("Type submitted is not supported")
	}
}
