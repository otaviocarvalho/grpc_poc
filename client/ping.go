package main

import (
	"log"
    "time"
    "fmt"
    "crypto/rand"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "grpc_poc/protobuf"
)

const (
	address = "localhost:50051"
)

func sendPing(client pb.PingPongClient, ping *pb.Ping) {
	_, err := client.SendPing(context.Background(), ping)
	if err != nil {
		log.Fatalf("Could not send message: %s", err)
	}

}

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Creates a new Ping
	client := pb.NewPingPongClient(conn)

    msgSize := 100
    randomBytes := make([]byte, msgSize)
    _, err = rand.Read(randomBytes)
    if err != nil {
        log.Fatal(err)
    }
    msg := string(randomBytes)

    ping := &pb.Ping{
        Payload: msg,
    }

    timeBucket := make(map[string]int64)

    numRuns := 100
    for i := 0; i < numRuns; i++ {
        clockBefore := time.Now()

        // Sends ping
        sendPing(client, ping)

        clockAfter := time.Now().Sub(clockBefore)
        timeBucket["sum"] += int64(clockAfter)
    }

    fmt.Printf("Total time\t%v\n", time.Duration(timeBucket["sum"]))
    fmt.Printf("Average latency\t%v\n", time.Duration(timeBucket["sum"] / int64(numRuns)))
}
