package main

import (
	"log"
    "time"
    //"math"
    //"fmt"
    "crypto/rand"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
    //"github.com/codahale/hdrhistogram"

	pb "grpc_poc/pingpong"
	hdr "github.com/otaviocarvalho/hdrhistogram"
)

const (
    //address = "191.232.175.110:50051"
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

    msgSize := 32 * 1000 // 64KB
    randomBytes := make([]byte, msgSize)
    _, err = rand.Read(randomBytes)
    if err != nil {
        log.Fatal(err)
    }
    msg := string(randomBytes)

    ping := &pb.Ping{
        Payload: msg,
    }

    // 1ms to 30 seconds range, 5 sigfigs precision
    hist := hdr.New(1000000, 30000000000, 5)

    numRuns := 10000
    startTime := time.Now()
    for i := 0; i < numRuns; i++ {
        startTimeLoop := time.Now()

        // Sends ping
        sendPing(client, ping)

        totalTimeLoop := time.Now().Sub(startTimeLoop)
        hist.RecordValue(totalTimeLoop.Nanoseconds())
    }

    totalTime := time.Now().Sub(startTime)
    log.Printf("QPS: %v", float64(numRuns) / totalTime.Seconds())

    log.Printf("50th: %d\n", hist.ValueAtQuantile(50))
    log.Printf("90th: %d\n", hist.ValueAtQuantile(90))
    log.Printf("99th: %d\n", hist.ValueAtQuantile(99))
    log.Printf("99.9th: %d\n", hist.ValueAtQuantile(99.9))
    log.Printf("99.99th: %d\n", hist.ValueAtQuantile(99.99))
    log.Printf("99.999th: %d\n", hist.ValueAtQuantile(99.999))
    log.Printf("99.9999th: %d\n", hist.ValueAtQuantile(99.9999))
    log.Printf("99.99999th: %d\n", hist.ValueAtQuantile(99.99999))
    log.Printf("99.999999th: %d\n", hist.ValueAtQuantile(99.999999))
    log.Printf("99.9999999th: %d\n", hist.ValueAtQuantile(99.9999999))
    log.Printf("99.99999999th: %d\n", hist.ValueAtQuantile(99.99999999))
    log.Printf("99.999999999th: %d\n", hist.ValueAtQuantile(99.999999999))
}
