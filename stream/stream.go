package stream

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Stream[T any] struct {
	Data chan T
}

func NewStream[T any](data []T) *Stream[T] {
	ch := make(chan T)
	go func() {
		for _, item := range data {
			ch <- item
		}
	}()
	return &Stream[T]{Data: ch}
}
func NewStreamFromPath(dataSetPath string) *Stream[string] {
	dataStream := NewStream(make([]string, 0))
	if dataSetPath == "" {
		go func() {
			for {
				random := rand.Intn(100)
				dataStream.Data <- strconv.Itoa(random)
			}
		}()
	} else {
		file, err := os.Open(dataSetPath)
		if err != nil {
			fmt.Println(err)
			panic("Could not open data set")
		}

		scanner := bufio.NewScanner(file)

		go func() {
			defer file.Close()
			for scanner.Scan() {
				line := scanner.Text()
				trimmedLine := strings.TrimSpace(line)
				if trimmedLine == "" {
					continue
				}
				dataStream.Data <- scanner.Text()
			}
		}()
	}
	return dataStream
}
