package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"sync"
	"sync/atomic"
	"time"

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

var latency = flag.Duration("l", 10*time.Millisecond, "artificial latency")

func (s *dataServer) SendMeasurement(ctx context.Context, in *pb.Measurement) (*pb.Measurement, error) {
	// Calculate EWMA
	expMAvgMutex.Lock()
	expMovingAvg.Add(in.GetValue())
	expMAvgMutex.Unlock()

	// Increase counter
	counterMutex.Lock()
	atomic.AddInt64(&counter, 1)
	curCount := counter
	counterMutex.Unlock()

	expMAvgMutex.RLock()
	measurement := &pb.Measurement{Id: curCount, Value: expMovingAvg.Value()}
	fmt.Println("counter:", curCount, "value:", expMovingAvg.Value())
	expMAvgMutex.RUnlock()

	return measurement, nil
}

func main() {
	flag.Parse()

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
