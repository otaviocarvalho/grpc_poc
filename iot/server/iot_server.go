package main

import (
	"log"
	"net"
	//"math"
    //"fmt"
    "sync/atomic"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "grpc_poc/iot"
)

const (
	port = ":50051"
)

type dataServer struct {}

var counter int64 = 0

func (s *dataServer) SendMeasurement(ctx context.Context, in *pb.Measurement) (*pb.Measurement, error) {
    //fmt.Println("Received ping: %s", counter)

    // Increase counter
    atomic.AddInt64(&counter, 1)

    return &pb.Measurement{Id: counter, Value: in.GetValue()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.MaxConcurrentStreams(1))
	pb.RegisterDataServer(s, &dataServer{})
	log.Printf("listening on port: %v", port)
	s.Serve(lis)
}
