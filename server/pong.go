package main

import (
	"log"
	"net"
    "crypto/rand"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "grpc_poc/protobuf"
)

const (
	port = ":50051"
)

type pingPongServer struct {}

func (s *pingPongServer) SendPing(ctx context.Context, in *pb.Ping) (*pb.Pong, error) {
	log.Printf("Received Ping with message: %s", in.Payload)
	log.Printf("Sending Pong back  with message: %s", "pong")

    msgSize := 100
    randomBytes := make([]byte, msgSize)
    _, err := rand.Read(randomBytes)
    if err != nil {
        log.Fatal(err)
    }
    msg := string(randomBytes)

    pong := &pb.Pong{
        Payload: msg,
    }

	return pong, nil
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
