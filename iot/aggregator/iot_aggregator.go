package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"go4.org/net/throttle"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/VividCortex/ewma"

	pb "grpc_poc/iot"
)

type dataServer struct{}

var counter int64 = 0

var expMovingAvg = ewma.NewMovingAverage()
var expMAvgMutex sync.RWMutex
var counterMutex sync.Mutex

var clientPort = flag.String("c", ":50051", "client ip/port")
var serverHost = flag.String("s", "localhost:50052", "ip/port")
var latency = flag.Duration("l", 10*time.Millisecond, "artificial latency")
var batchSize = flag.Int64("b", 1, "batch size")

var messageChannel = make(chan int64)

func (s *dataServer) SendMeasurement(ctx context.Context, in *pb.Measurement) (*pb.Measurement, error) {
	// Calculate EWMA
	expMAvgMutex.Lock()
	expMovingAvg.Add(in.GetValue())
	expMAvgMutex.Unlock()

	// Increase counter
	counterMutex.Lock()
	atomic.AddInt64(&counter, 1)
	messageChannel <- counter
	counterMutex.Unlock()

	expMAvgMutex.RLock()
	measurement := &pb.Measurement{Id: counter, Value: expMovingAvg.Value()}
	expMAvgMutex.RUnlock()

	return measurement, nil
}

func main() {
	flag.Parse()

	exitChannel := make(chan struct{})

	// Receive connections from client
	go func() {
		lis, err := net.Listen("tcp", *clientPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		rate := throttle.Rate{Latency: *latency}
		lis = &throttle.Listener{lis, rate, rate}

		s := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32))
		pb.RegisterDataServer(s, &dataServer{})
		log.Printf("listening on port: %v", *clientPort)
		s.Serve(lis)
	}()

	// Sends aggregated data to server
	go func() {
		// Set up a connection to the gRPC server.
		conn, err := grpc.Dial(*serverHost, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		// Creates a new data client
		client := pb.NewDataClient(conn)

		// Transmits aggregated messages from client to global server
		for {
			counter := <-messageChannel

			if counter%*batchSize == 0 {

				data := &pb.Measurement{
					Value: expMovingAvg.Value(),
				}

				_, err = client.SendMeasurement(context.Background(), data)
				if err != nil {
					log.Fatalf("Could not send message: %s", err)
				}

				fmt.Println(counter)
			}
		}
	}()

	<-exitChannel
}
