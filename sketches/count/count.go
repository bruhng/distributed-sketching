package count

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"slices"

	"github.com/spaolacci/murmur3"
)

type CountSketch[T any, R any] struct {
	Sketch [][]int
	Seeds  []uint32
}

func NewCountSketch[T any](seed int64, size uint64, num_hashes int) *CountSketch[T, int] {
	arr := make([][]int, num_hashes)

	for i := 0; i < num_hashes; i++ {
		arr[i] = make([]int, size)
	}

	r := rand.New(rand.NewSource(seed))
	seeds := make([]uint32, num_hashes)
	for i := 0; i < num_hashes; i++ {
		seeds[i] = r.Uint32()
	}
	return &CountSketch[T, int]{Sketch: arr, Seeds: seeds}
}

func getSign(data []byte) int {
	hash := murmur3.New64()
	hash.Write(data)
	hashValue := hash.Sum64()

	if hashValue%2 == 0 {
		return 1
	}
	return -1
}

func getIndex(data []byte, seed uint32, size uint64) uint64 {
	hash := murmur3.New64WithSeed(seed)
	hash.Write(data)
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
	fmt.Println(cs)
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
		return (results[result_size/2-1] + results[result_size/2]) / 2
	}
	return results[result_size/2]
}

func (cs *CountSketch[T, R]) Merge(sketch CountSketch[T, R]) {
	/*if reflect.DeepEqual(cs.Seeds, sketch.Seeds) {
		panic("Missmatched hash function in merged sketches")
	}*/
	if len(cs.Sketch[0]) != len(sketch.Sketch[0]) {
		panic("Missmatched length of second dimension in merged sketches")
	}
	for i, rows := range cs.Sketch {
		for j, elems := range rows {
			cs.Sketch[i][j] = elems + sketch.Sketch[i][j]
		}
	}
}
