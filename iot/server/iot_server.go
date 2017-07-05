package main

import (
    "log"
    "net"
    "math"
    "flag"

    "golang.org/x/net/context"
    "google.golang.org/grpc"

    pb "grpc_poc/iot"
)

type dataServer struct {}

var port = flag.String("p", ":50051", "ip/port")

func (s *dataServer) SendMeasurement(ctx context.Context, in *pb.Measurement) (*pb.Measurement, error) {
    return &pb.Measurement{Value: in.GetValue()}, nil
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
