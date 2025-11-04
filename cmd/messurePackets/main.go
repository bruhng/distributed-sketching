package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bruhng/distributed-sketching/client"
	pb "github.com/bruhng/distributed-sketching/proto"
	"github.com/bruhng/distributed-sketching/shared"
	"github.com/bruhng/distributed-sketching/stream"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var SERVER_ADR = "127.0.0.1"
var PORT = "8080"
var clientAmmount = 100
var streamRate = 100000
var mergeRate = 1000
var NUM_MERGES = 10000
var DATA_SET_PATH = "../../data/PVS 1/dataset_gps.csv"
var HEADER_NAME = "speed_meters_per_second"
var meregesMade uint64 = 0

func startRealConnection(adr string) (pb.SketcherClient, *grpc.ClientConn, error) {

	conn, err := grpc.NewClient(adr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	c := pb.NewSketcherClient(conn)
	conn.Connect()
	// wait for the connection to be established
	for {
		state := conn.GetState()
		if state == connectivity.Idle || state == connectivity.Connecting {
			continue
		} else if state == connectivity.Ready {
			break
		} else {
			return c, conn, fmt.Errorf("Could not establish connection to server")
		}
	}
	//
	return c, conn, nil
}

func startKll[T shared.Number](wg *sync.WaitGroup, fg *sync.WaitGroup, cond *sync.Cond, streamRate int, mergeRate int, sketch *pb.KLLSketch, merges int) {
	c, conn, err := startRealConnection(SERVER_ADR + ":" + PORT)

	var reconAttempt *int = new(int)
	if err != nil {
		panic("test failed because no connection")
	}
	wg.Done()
	cond.L.Lock()

	cond.Wait()
	cond.L.Unlock()

	for range merges {

		time.Sleep(time.Duration(streamRate*mergeRate) * time.Nanosecond)
		go func() {
			client.MakeRequest[pb.KLLSketch](sketch, SERVER_ADR, c.MergeKll, conn, &c, startRealConnection, reconAttempt)
			atomic.AddUint64(&meregesMade, 1)
		}()
	}
	return
}

func main() {

	var wg sync.WaitGroup
	var fg sync.WaitGroup
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	dataStream := *stream.NewStreamFromCsv[float64](DATA_SET_PATH, HEADER_NAME, streamRate, -1)
	sketch := client.GetKll(100, mergeRate, dataStream)

	for range clientAmmount {
		wg.Add(1)

		go startKll[float64](&wg, &fg, cond, streamRate, mergeRate, sketch, NUM_MERGES)
	}
	wg.Wait()

	cond.Broadcast()
	optimal := 5.0 / (float64(streamRate) * float64(mergeRate)) * float64(clientAmmount) * 1e9
	for {
		time.Sleep(time.Duration(5) * time.Second)
		fmt.Println(atomic.SwapUint64(&meregesMade, 0), optimal)
	}

}
