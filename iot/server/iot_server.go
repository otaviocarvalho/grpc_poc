package main

import (
    "log"
    "net"
    "math"
    "flag"
    "sync/atomic"
    "sync"

    "golang.org/x/net/context"
    "google.golang.org/grpc"

    "github.com/VividCortex/ewma"

    pb "grpc_poc/iot"
)

type dataServer struct {}

var counter int64 = 0

var expMovingAvg = ewma.NewMovingAverage()

var expMAvgMutex sync.RWMutex

var counterMutex sync.Mutex

var port = flag.String("p", ":50051", "ip/port")

func (s *dataServer) SendMeasurement(ctx context.Context, in *pb.Measurement) (*pb.Measurement, error) {
    // Calculate EWMA
    expMAvgMutex.Lock()
    expMovingAvg.Add(in.GetValue())
    expMAvgMutex.Unlock()

    // Increase counter
    counterMutex.Lock()
    atomic.AddInt64(&counter, 1)
    counterMutex.Unlock()

    expMAvgMutex.RLock()
    measurement := &pb.Measurement{Id: counter, Value: expMovingAvg.Value()}
    expMAvgMutex.RUnlock()

    return measurement, nil
}

func main() {
    flag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32))
	pb.RegisterDataServer(s, &dataServer{})
	log.Printf("listening on port: %v", *port)
	s.Serve(lis)
}
