package main

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"net"
	"strings"
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

var clientPort = flag.String("cp", ":50051", "client ip/port")
var serverPort = flag.String("sp", ":50052", "server ip/port")
var latency = flag.Duration("l", 10*time.Millisecond, "artificial latency")

func (s *dataServer) SendMeasurement(ctx context.Context, in *pb.Measurement) (*pb.Measurement, error) {
	// Calculate EWMA
	expMAvgMutex.Lock()
	expMovingAvg.Add(in.GetValue())
	expMAvgMutex.Unlock()

	// Increase counter
	counterMutex.Lock()
	atomic.AddInt64(&counter, 1)
	counterMutex.Unlock()

	expMAvgMutex.RLock()
	measurement := &pb.Measurement{Id: counter, Value: expMovingAvg.Value()}
	expMAvgMutex.RUnlock()

	return measurement, nil
}

func sendMeasurementToGlobalServer(client pb.DataClient, data *pb.Measurement) {

	_, err := client.SendMeasurement(context.Background(), data)
	if err != nil {
		log.Fatalf("Could not send message: %s", err)
	}

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
		conn, err := grpc.Dial(strings.Join([]string{"localhost", *serverPort}, ""), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		// Creates a new data client
		client := pb.NewDataClient(conn)

		data := &pb.Measurement{
			Value: rand.ExpFloat64(),
		}

		_, err = client.SendMeasurement(context.Background(), data)
		if err != nil {
			log.Fatalf("Could not send message: %s", err)
		}

	}()

	<-exitChannel
}
