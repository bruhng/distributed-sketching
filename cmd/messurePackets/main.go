package main

import (
	"flag"
	"fmt"
	"math"
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
	"google.golang.org/protobuf/proto"
)

var SERVER_ADR = "127.0.0.1"
var PORT = "8080"
var NUM_MERGES = 10000
var DATA_SET_PATH = "../../data/PVS 1/dataset_gps.csv"
var HEADER_NAME = "speed_meters_per_second"
var meregesMade uint64 = 0
var mesuringInterval = 5.0
var timeUnit = time.Nanosecond
var timeExponent = 1 / (float64(timeUnit) / float64(time.Second))

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

		time.Sleep(time.Duration(streamRate*mergeRate) * timeUnit)
		go func() {
			client.MakeRequest[pb.KLLSketch](sketch, SERVER_ADR, c.MergeKll, conn, &c, startRealConnection, reconAttempt)
			atomic.AddUint64(&meregesMade, 1)
		}()
	}
	return
}
func startCount[T shared.Number](wg *sync.WaitGroup, fg *sync.WaitGroup, cond *sync.Cond, streamRate int, mergeRate int, sketch *pb.CountSketch, merges int) {
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

		time.Sleep(time.Duration(streamRate*mergeRate) * timeUnit)
		go func() {
			client.MakeRequest[pb.CountSketch](sketch, SERVER_ADR, c.MergeCount, conn, &c, startRealConnection, reconAttempt)
			atomic.AddUint64(&meregesMade, 1)
		}()
	}
	return
}
func startBadKll[T shared.Number](wg *sync.WaitGroup, fg *sync.WaitGroup, cond *sync.Cond, streamRate int, mergeRate int, sketch *pb.BadArray, merges int) {
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

		time.Sleep(time.Duration(streamRate*mergeRate) * timeUnit)
		go func() {
			client.MakeRequest[pb.BadArray](sketch, SERVER_ADR, c.BadKll, conn, &c, startRealConnection, reconAttempt)
			atomic.AddUint64(&meregesMade, 1)
		}()
	}
	return
}
func startBadCount[T shared.Number](wg *sync.WaitGroup, fg *sync.WaitGroup, cond *sync.Cond, streamRate int, mergeRate int, sketch *pb.BadArray, merges int) {
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

		time.Sleep(time.Duration(streamRate*mergeRate) * timeUnit)
		go func() {
			client.MakeRequest[pb.BadArray](sketch, SERVER_ADR, c.BadCount, conn, &c, startRealConnection, reconAttempt)
			atomic.AddUint64(&meregesMade, 1)
		}()
	}
	return
}

func main() {
	var wg sync.WaitGroup
	var fg sync.WaitGroup
	var mu sync.Mutex
	mergeRate := *flag.Int("mergeRate", 1000, "merge rate for clients")
	clientAmount := *flag.Int("clientAmount", 1, "number of clients")
	streamRate := *flag.Int("streamRate", 1000, "stream rate for clients")
	sketchType := *flag.String("type", "kll", "sketch type")

	cond := sync.NewCond(&mu)
	dataStream := *stream.NewStreamFromCsv[float64](DATA_SET_PATH, HEADER_NAME, streamRate, -1)
	kllSketch := client.GetKll(100, mergeRate, dataStream)
	countSketch := client.GetCount(mergeRate, dataStream)
	badArr := client.GetBad(mergeRate, dataStream)
	var sketchSize uint64

	for range clientAmount {
		wg.Add(1)

		if sketchType == "kll" {
			go startKll[float64](&wg, &fg, cond, streamRate, mergeRate, kllSketch, NUM_MERGES)
		}
		if sketchType == "count" {
			go startCount[float64](&wg, &fg, cond, streamRate, mergeRate, countSketch, NUM_MERGES)
		}
		if sketchType == "badCount" {
			go startBadCount[float64](&wg, &fg, cond, streamRate, mergeRate, badArr, NUM_MERGES)
		}
		if sketchType == "badKll" {
			go startBadKll[float64](&wg, &fg, cond, streamRate, mergeRate, badArr, NUM_MERGES)
		}

	}
	if sketchType == "kll" {
		sketchSize = uint64(proto.Size(kllSketch))
	}
	if sketchType == "count" {
		sketchSize = uint64(proto.Size(kllSketch))
	}
	if sketchType == "badCount" {
		sketchSize = uint64(proto.Size(kllSketch))
	}
	if sketchType == "badKll" {
		sketchSize = uint64(proto.Size(kllSketch))
	}
	wg.Wait()

	cond.Broadcast()
	optimal := math.Round(mesuringInterval / (float64(streamRate) * float64(mergeRate)) * float64(clientAmmount) * timeExponent)
	optBandWidth := int64(math.Min(1e9/3.0, float64(sketchSize)*optimal/mesuringInterval))

	for {
		time.Sleep(time.Duration(mesuringInterval) * time.Second)
		merges := atomic.SwapUint64(&meregesMade, 0)
		fmt.Println(merges, optimal, sketchSize*merges/uint64(mesuringInterval), optBandWidth)
	}

}
