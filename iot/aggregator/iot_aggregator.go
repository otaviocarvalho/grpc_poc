package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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

	hdr "github.com/otaviocarvalho/hdrhistogram"
)

type dataServer struct{}

var counter int64 = 0

var expMovingAvg = ewma.NewMovingAverage()
var expMAvgMutex sync.RWMutex
var counterMutex sync.Mutex

var clientPort = flag.String("c", ":50051", "client ip/port")
var serverHost = flag.String("s", "localhost:50052", "ip/port")
var latency = flag.Duration("l", 0*time.Millisecond, "artificial latency")
var batchSize = flag.Int64("b", 1, "batch size")
var batchLogSize = flag.Int64("blog", 1000, "batch size for logging")
var outputFile = flag.String("o", "./output_stats_aggregator.json", "output json file")
var enableLogs = flag.Bool("log", false, "enable/disable logs")

var messageChannel = make(chan int64, 100)

type Stats struct {
	Perc50    int64 `json:"p50"`
	Perc90    int64 `json:"p90"`
	Perc99    int64 `json:"p99"`
	Perc999   int64 `json:"p999"`
	Perc9999  int64 `json:"p9999"`
	Perc99999 int64 `json:"p99999"`
}

func (s *dataServer) SendMeasurement(ctx context.Context, in *pb.Measurement) (*pb.Measurement, error) {
	// Calculate EWMA
	expMAvgMutex.Lock()
	expMovingAvg.Add(in.GetValue())
	expMAvgMutex.Unlock()

	// Increase counter
	counterMutex.Lock()
	atomic.AddInt64(&counter, 1)
	curCount := counter
	messageChannel <- curCount
	counterMutex.Unlock()

	expMAvgMutex.RLock()
	measurement := &pb.Measurement{Id: curCount, Value: expMovingAvg.Value()}
	expMAvgMutex.RUnlock()

	return measurement, nil
}

func plotStats(numRuns int64, hist *hdr.Histogram) {
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

func saveStats(numRuns int64, hist *hdr.Histogram) {
	stats := Stats{
		hist.ValueAtQuantile(50),
		hist.ValueAtQuantile(90),
		hist.ValueAtQuantile(99),
		hist.ValueAtQuantile(999),
		hist.ValueAtQuantile(9999),
		hist.ValueAtQuantile(99999),
	}

	statsJson, _ := json.Marshal(stats)
	err := ioutil.WriteFile(*outputFile, statsJson, 0664)
	if err != nil {
		fmt.Println("Writing output file", err.Error())
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
		conn, err := grpc.Dial(*serverHost, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		// Creates a new data client
		client := pb.NewDataClient(conn)

		// Transmits aggregated messages from client to global server
		hist := hdr.New(1000000, 30000000000, 5)
		var histMutex sync.RWMutex
		var counterAux = int64(0)
		for {
			var isValidMeasurement = true

			select {
			case c := <-messageChannel:
				counterAux = c
			default:
				isValidMeasurement = false
			}

			startTimeLoop := time.Now()

			if counterAux%*batchSize == 0 && isValidMeasurement {

				data := &pb.Measurement{
					Value: expMovingAvg.Value(),
				}

				_, err = client.SendMeasurement(context.Background(), data)
				if err != nil {
					log.Fatalf("Could not send message: %s", err)
				}

				// Save histogram for each request
				totalTimeLoop := time.Now().Sub(startTimeLoop)
				histMutex.Lock()
				hist.RecordValue(totalTimeLoop.Nanoseconds())
				histMutex.Unlock()

				// Write histogram for a batch of requests
				if (*enableLogs) && (counterAux%*batchLogSize == 0) {
					plotStats(int64(*batchSize), hist)
					saveStats(int64(*batchSize), hist)
				}
			}
		}
	}()

	<-exitChannel
}
