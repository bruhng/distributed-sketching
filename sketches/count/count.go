package count

import (
	"hash"
	"hash/fnv"
	"math/rand"
	"strconv"

	sk "github.com/bruhng/distributed-sketching/sketches"
)

type CountSketch[T any] struct {
	sketch [][]int
	hashes []hash.Hash64
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

}

func getIndex(data []byte, seed int, size uint64) uint64 {
	seededBytes := []byte(strconv.Itoa(seed))
	seededBytes = append(seededBytes, data...)
	hash := fnv.New64()
	hash.Write(seededBytes)
	hashVal := hash.Sum64()
	return hashVal % size

}

func (cs *CountSketch[T]) Add(item T) {
}

func (cs *CountSketch[T]) Query(item T) int {
}

func (cs *CountSketch[T]) Merge(sketch CountSketch[T]) {
}
