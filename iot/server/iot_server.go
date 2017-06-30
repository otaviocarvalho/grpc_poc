package main

import (
    "log"
    "net"
    "math"
    "flag"
    "sync/atomic"

    "golang.org/x/net/context"
    "google.golang.org/grpc"

    "github.com/VividCortex/ewma"

    pb "grpc_poc/iot"
)

type dataServer struct {}

var counter int64 = 0

var expMovingAvg = ewma.NewMovingAverage()

var port = flag.String("p", ":50051", "ip/port")

func (s *dataServer) SendMeasurement(ctx context.Context, in *pb.Measurement) (*pb.Measurement, error) {
    // Calculate EWMA
    expMovingAvg.Add(in.GetValue())

    // Increase counter
    atomic.AddInt64(&counter, 1)

    return &pb.Measurement{Id: counter, Value: expMovingAvg.Value()}, nil
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
