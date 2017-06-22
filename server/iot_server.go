package main

import (
	"log"
	"net"
    //"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "grpc_poc/protobuf_iot"
)

const (
	port = ":50051"
)

type dataServer struct {}

func (s *dataServer) SendMeasurement(ctx context.Context, in *pb.Measurement) (*pb.Measurement, error) {
    //fmt.Println("Received ping with message: %s", in.Payload)
    return &pb.Measurement{Value: in.GetValue()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDataServer(s, &dataServer{})
	log.Printf("listening on port: %v", port)
	s.Serve(lis)
}
