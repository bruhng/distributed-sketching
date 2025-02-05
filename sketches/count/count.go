package count

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	"math/rand"
	"reflect"
	"slices"
	"strconv"

	sk "github.com/bruhng/distributed-sketching/sketches"
)

type CountSketch[T any, R any] struct {
	Sketch [][]int
	Seeds  []int
}

func NewCountSketch[T any](seed int64, size uint64, num_hashes int) sk.Sketch[T, int] {
	arr := make([][]int, num_hashes)

	for i := 0; i < num_hashes; i++ {
		arr[i] = make([]int, size)
	}

	rand.Seed(seed)
	seeds := make([]int, num_hashes)
	for i := 0; i < num_hashes; i++ {
		seeds[i] = rand.Intn(2 ^ 63)
	}

	return CountSketch[T, int]{Sketch: arr, Seeds: seeds}
}

func getSign(data []byte) int {
	hash := fnv.New32a()
	hash.Write(data)
	hashValue := hash.Sum32()

	if hashValue%2 == 0 {
		return 1
	}
	return -1
}

func getIndex(data []byte, seed int, size uint64) uint64 {
	seededBytes := []byte(strconv.Itoa(seed))
	seededBytes = append(seededBytes, data...)
	hash := fnv.New64()
	hash.Write(seededBytes)
	hashVal := hash.Sum64()
	return hashVal % size

}

func (cs *CountSketch[T, R]) Add(item T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(item)
	if err != nil {
		panic("hi i am not handled (could not convert data to bytes)")
	}
	size := uint64(len(cs.Sketch[0]))
	sign := getSign(buf.Bytes())
	for i, seed := range cs.Seeds {
		index := getIndex(buf.Bytes(), seed, size)
		cs.Sketch[i][index] += sign
	}
}

func (cs *CountSketch[T, R]) Query(item T) int {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(item)
	if err != nil {
		panic("hi i am not handled (could not convert data to bytes)")
	}
	size := uint64(len(cs.Sketch[0]))
	sign := getSign(buf.Bytes())
	result_size := len(cs.Seeds)
	results := make([]int, result_size)
	for i, seed := range cs.Seeds {
		index := getIndex(buf.Bytes(), seed, size)
		results[i] = cs.Sketch[i][index] * sign
	}

	slices.Sort(results)
	if result_size%2 == 0 {
		return (results[result_size] + results[result_size/2+1]) / 2
	}
	return results[result_size/2]
}

func (cs *CountSketch[T, R]) Merge(sketch CountSketch[T, int]) {
	if reflect.DeepEqual(cs.Seeds, sketch.Seeds) {
		panic("Missmatched hash function in merged sketches")
	}
	if len(cs.Sketch[0]) != len(sketch.Sketch[0]) {
		panic("Missmatched length of second dimension in merged sketches")
	}
	for i, rows := range cs.Sketch {
		for j, elems := range rows {
			cs.Sketch[i][j] = elems + sketch.Sketch[i][j]
		}
	}
}
