package main

import (
	"log"
	"net"
    "fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "grpc_poc/protobuf"
)

const (
	port = ":50051"
)

type pingPongServer struct {}

func (s *pingPongServer) SendPing(ctx context.Context, in *pb.Ping) (*pb.Pong, error) {
    fmt.Println("Received ping with message: %s", in.Payload)
    return &pb.Pong{ Payload: in.Payload }, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPingPongServer(s, &pingPongServer{})
	log.Printf("listening on port: %v", port)
	s.Serve(lis)
}
