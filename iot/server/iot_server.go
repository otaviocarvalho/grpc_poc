package main

import (
	"fmt"
	"flag"
	"log"
	"math"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"encoding/json"
	"io/ioutil"

	"go4.org/net/throttle"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/VividCortex/ewma"

	pb "grpc_poc/iot"
)

type dataServer struct{}

var counter int64 = 0

var expMovingAvg = ewma.NewMovingAverage()

var expMAvgMutex sync.RWMutex

var counterMutex sync.Mutex

var port = flag.String("p", ":50052", "ip/port")
var latency = flag.Duration("l", 0*time.Millisecond, "artificial latency")
var outputFile = flag.String("o", "./output_stats_server.json", "output json file")
var batchLogSize = flag.Int64("blog", 1000, "batch size for logging")

var messageChannel = make(chan int64, 100000000)

var startTime time.Time

func (s *dataServer) SendMeasurement(ctx context.Context, in *pb.Measurement) (*pb.Measurement, error) {
	// Calculate EWMA
	expMAvgMutex.Lock()
	expMovingAvg.Add(in.GetValue())
	expMAvgMutex.Unlock()

	// Increase counter
	counterMutex.Lock()
	atomic.AddInt64(&counter, 1)
	curCount := counter
	messageChannel <- curCount
	counterMutex.Unlock()

	expMAvgMutex.RLock()
	measurement := &pb.Measurement{Id: curCount, Value: expMovingAvg.Value()}
	expMAvgMutex.RUnlock()

	return measurement, nil
}

type stats struct {
	QPS	float64 `json:"qps"` // Should multiply this number by the batch size on the aggregator side
}

func saveStats(numRuns int64, totalTime time.Duration) {
	stats := stats{
		float64(numRuns) / totalTime.Seconds(),
	}

	statsJSON, _ := json.Marshal(stats)
	err := ioutil.WriteFile(*outputFile, statsJSON, 0664)
	if err != nil {
		fmt.Println("Writing output file", err.Error())
	}
}

func main() {
	flag.Parse()

	// Log summary each couple of seconds
	startTime = time.Now()
	go func() {
		for {
			select {
				case curCount := <-messageChannel:
					fmt.Println(curCount)
					if (curCount % *batchLogSize == 0) {
						saveStats(curCount, time.Now().Sub(startTime))
					}
			}
		}	
	}()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rate := throttle.Rate{Latency: *latency}
	lis = &throttle.Listener{lis, rate, rate}

	s := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32))
	pb.RegisterDataServer(s, &dataServer{})
	log.Printf("listening on port: %v", *port)
	s.Serve(lis)
}
