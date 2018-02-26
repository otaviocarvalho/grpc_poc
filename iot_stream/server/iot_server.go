package main

import (
	"flag"
	"io"
	"log"
	"math"
	"net"
	"sync/atomic"

	"google.golang.org/grpc"

	pb "grpc_poc/iot_stream"
)

type MeterCommunicatorServer struct{}

var counter int64 = 0

var port = flag.String("p", ":50051", "ip/port")

func (s *MeterCommunicatorServer) SimpleRPC(stream pb.MeterCommunicator_SimpleRPCServer) error {

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			log.Fatalf("error (EOF) receiving msg: %v", err)
			return nil
		}
		if err != nil {
			log.Fatalf("error receiving msg: %v", err)
			return err
		}

		// Increase counter
		atomic.AddInt64(&counter, 1)
	}
}

func main() {
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32))
	pb.RegisterMeterCommunicatorServer(s, &MeterCommunicatorServer{})
	log.Printf("listening on port: %v", *port)
	s.Serve(lis)
}
