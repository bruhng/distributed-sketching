package stream

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/bruhng/distributed-sketching/shared"
)

type Stream[T shared.Number] struct {
	Data chan T
}

func NewStream[T shared.Number](data []T) *Stream[T] {
	ch := make(chan T)
	go func() {
		for _, item := range data {
			ch <- item
		}
	}()
	return &Stream[T]{Data: ch}
}

// func NewStreamFromPath(dataSetPath string) *Stream[string] {
// 	dataStream := NewStream(make([]string, 0))
// 	if dataSetPath == "" {
// 		go func() {
// 			for {
// 				random := rand.Intn(100)
// 				dataStream.Data <- strconv.Itoa(random)
// 			}
// 		}()
// 	} else {
// 		file, err := os.Open(dataSetPath)
// 		if err != nil {
// 			fmt.Println(err)
// 			panic("Could not open data set")
// 		}

// 		scanner := bufio.NewScanner(file)

// 		go func() {
// 			defer file.Close()
// 			for scanner.Scan() {
// 				line := scanner.Text()
// 				trimmedLine := strings.TrimSpace(line)
// 				if trimmedLine == "" {
// 					continue
// 				}
// 				dataStream.Data <- scanner.Text()
// 			}
// 		}()
// 	}
// 	return dataStream
// }

func NewStreamFromCsv[T shared.Number](csvPath string, field string) *Stream[T] {
	dataStream := NewStream(make([]T, 0))
	file, err := os.Open(csvPath)
	if err != nil {
		panic("Could not read csv")
	}

	go func() {

		defer file.Close()

		reader := csv.NewReader(file)

		header, err := reader.Read()
		if err != nil {
			panic("could not reader header")
		}

		columnIndex := -1
		for i, h := range header {
			if h == field {
				columnIndex = i
				break
			}
		}
		fmt.Println("columnIndex: ", columnIndex)
		if columnIndex == -1 {
			panic("Invalid field name")
		}

		for {
			record, err := reader.Read()
			if err != nil {
				break
			}
			data := record[columnIndex]
			parsedData, err := parseNumber(data)
			if err != nil {
				panic("Data is not int or float")
			}
			parsed, ok := parsedData.(T)
			if ok {
				dataStream.Data <- parsed
			} else {
				panic("Data is not int or float")
			}

		}
		fmt.Println("stream is now empty")

	}()
	return dataStream
}

func parseNumber(s string) (any, error) {
	if intValue, err := strconv.ParseInt(s, 10, 64); err == nil {
		return intValue, nil
	}

	if floatValue, err := strconv.ParseFloat(s, 64); err == nil {
		return floatValue, nil
	}

	return nil, fmt.Errorf("%s is not a valid number", s)
}
