package kll

import (
	"cmp"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"sort"
)

type KLLSketch[T cmp.Ordered, R any] struct {
	Sketch [][]T
	K      int
	N      int
}

func NewKLLSketch[T cmp.Ordered, R int](k int) *KLLSketch[T, int] {
	arr := make([][]T, 1)
	return &KLLSketch[T, int]{Sketch: arr, K: k}
}

func getSize(k int, h int, H int) int {
	return max(2, k*int(math.Round(math.Pow(2.0/3.0, float64(H-h)))))
}

func (kll *KLLSketch[T, R]) Add(item T) {
	kll.Sketch[0] = append(kll.Sketch[0], item)
	kll.N++
	compress(kll)
}

func everyOther[T any](xs []T) []T {
	if len(xs) == 0 {
		return []T{}
	}
	if len(xs) == 1 {
		return xs
	}
	return append([]T{xs[0]}, everyOther(xs[2:])...)
}

func compress[T cmp.Ordered, R any](kll *KLLSketch[T, R]) {
	h := 0
	for {
		row := kll.Sketch[h]
		if len(row) >= getSize(kll.K, h, len(kll.Sketch)) {
			if len(kll.Sketch) == h+1 {
				kll.Sketch = append(kll.Sketch, make([]T, 0))
			}
			slices.Sort(row)
			even := rand.Int() % 2
			if even == 0 {
				kll.Sketch[h+1] = append(kll.Sketch[h+1], everyOther(row)...)
			} else {
				kll.Sketch[h+1] = append(kll.Sketch[h+1], everyOther(row[1:])...)
			}
			kll.Sketch[h] = make([]T, 0)
		}
		if len(kll.Sketch) == h+1 {
			return
		}
		h++

	}
}

func (kll *KLLSketch[T, int]) Merge(sketch KLLSketch[T, int]) {
	H := max(len(kll.Sketch), len(sketch.Sketch))
	diff := H - len(kll.Sketch)
	kll.Sketch = append(kll.Sketch, make([][]T, diff)...)
	for h := 0; h < len(sketch.Sketch); h++ {
		kll.Sketch[h] = append(kll.Sketch[h], sketch.Sketch[h]...)
	}

	kll.N += sketch.N
	compress(kll)
}

func (kll *KLLSketch[T, R]) Query(val T) int {
	sum := 0
	for h, row := range kll.Sketch {
		for _, elem := range row {
			if elem <= val {
				sum += int(math.Pow(2.0, float64(h)))
			}
		}
	}
	return sum
}

func (kll *KLLSketch[T, R]) QueryQuantile(q int) (T, error) {
	if q > kll.N {
		return kll.Sketch[0][0], fmt.Errorf("Requested quantile is bigger than N")
	}
	quantileSum := 0
	sketch := kll.Sketch
	var smallestH int
	var smallestVal T
	for _, row := range sketch {
		sort.Slice(row, func(i, j int) bool {
			return row[i] < row[j]
		})
	}

	for quantileSum < q {
		smallestH = -1
		for h, row := range sketch {
			if len(row) == 0 {
				continue
			}
			if smallestH == -1 {
				smallestVal = row[0]
				smallestH = h
			}
			if row[0] < smallestVal {
				smallestVal = row[0]
				smallestH = h
			}
		}
		sketch[smallestH] = sketch[smallestH][1:]
		quantileSum += int(math.Pow(2.0, float64(smallestH)))
	}
	return smallestVal, nil
}
