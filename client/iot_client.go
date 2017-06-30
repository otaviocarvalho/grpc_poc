package main

import (
	"log"
    "time"
    //"math"
    //"fmt"
    "math/rand"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
    //"github.com/codahale/hdrhistogram"

	pb "grpc_poc/protobuf_iot"
	hdr "github.com/otaviocarvalho/hdrhistogram"
)

const (
    //address = "191.232.175.110:50051"
    address = "localhost:50051"
)

func sendMeasurement(client pb.DataClient, data *pb.Measurement) {

    _, err := client.SendMeasurement(context.Background(), data)
	if err != nil {
		log.Fatalf("Could not send message: %s", err)
	}

    //fmt.Println("Measurement response id: %v", measurementResponse.GetId())
}

func plotQuantiles(numRuns int, totalTime time.Duration, hist *hdr.Histogram) {
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

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Creates a new data client
	client := pb.NewDataClient(conn)

    data := &pb.Measurement{
        Value: rand.ExpFloat64(),
    }

    // 1ms to 30 seconds range, 5 sigfigs precision
    hist := hdr.New(1000000, 30000000000, 5)

    numRuns := 10000
    startTime := time.Now()
    for i := 0; i < numRuns; i++ {
        startTimeLoop := time.Now()

        // Sends measurement
        sendMeasurement(client, data)

        totalTimeLoop := time.Now().Sub(startTimeLoop)
        hist.RecordValue(totalTimeLoop.Nanoseconds())
    }

    totalTime := time.Now().Sub(startTime)

    plotQuantiles(numRuns, totalTime, hist)
}
