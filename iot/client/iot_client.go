package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"github.com/golang/protobuf/proto"
	"github.com/VividCortex/ewma"

	pb "grpc_poc/iot"

	hdr "github.com/otaviocarvalho/hdrhistogram"
)

var concurrency = flag.Int("c", 1, "concurrency")
var batchSize = flag.Int("b", 1, "batch size")
var total = flag.Int("n", 1, "total requests for all clients")
var host = flag.String("s", "localhost:50051", "ip/port")
var outputFile = flag.String("o", "./output_stats.json", "output json file")

var expMovingAvg = ewma.NewMovingAverage()

var expMAvgMutex sync.RWMutex

func sendMeasurement(client pb.DataClient, data *pb.Measurement) {

	_, err := client.SendMeasurement(context.Background(), data)
	if err != nil {
		log.Fatalf("Could not send message: %s", err)
	}

}

func processMeasurement(client pb.DataClient, data *pb.Measurement) float64 {
	// Calculate EWMA
	expMAvgMutex.Lock()
	expMovingAvg.Add(data.GetValue())
	expMAvgMutex.Unlock()

	expMAvgMutex.RLock()
	measurementValue := expMovingAvg.Value()
	expMAvgMutex.RUnlock()

	return measurementValue
}

func plotStats(numRuns int, totalTime time.Duration, hist *hdr.Histogram) {
	log.Printf("QPS: %v", float64(numRuns)/totalTime.Seconds())

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

type Stats struct {
	Qps       float64 `json:"qps"`
	Perc50    int64   `json:"p50"`
	Perc90    int64   `json:"p90"`
	Perc99    int64   `json:"p99"`
	Perc999   int64   `json:"p999"`
	Perc9999  int64   `json:"p9999"`
	Perc99999 int64   `json:"p99999"`
}

func saveStats(numRuns int, totalTime time.Duration, hist *hdr.Histogram) {
	stats := Stats{
		float64(numRuns) / totalTime.Seconds(),
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

	//fmt.Printf("%+v", stats)
}

func main() {
	flag.Parse()
	n := *concurrency
	m := *total / n

	// 1ms to 30 seconds range, 5 sigfigs precision
	hist := hdr.New(1000000, 30000000000, 5)

	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(*host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// Creates a new data client
	client := pb.NewDataClient(conn)

	data := &pb.Measurement{
		Value: rand.ExpFloat64(),
	}

	var wg sync.WaitGroup
	wg.Add(*total)

	var histMutex sync.RWMutex

	startTime := time.Now()
	for i := 0; i < n; i++ {
		go func() {

			for j := 0; j < m; j++ {
				startTimeLoop := time.Now()

				// Controls processing local and remote
				if (j+1)%*batchSize != 0 {
					processMeasurement(client, data)
				} else {
					sendMeasurement(client, data)
				}

				totalTimeLoop := time.Now().Sub(startTimeLoop)

				histMutex.Lock()
				hist.RecordValue(totalTimeLoop.Nanoseconds())
				histMutex.Unlock()

				wg.Done()
			}

		}()
	}

	wg.Wait()

	totalTime := time.Now().Sub(startTime)

	plotStats(*total, totalTime, hist)
	saveStats(*total, totalTime, hist)

	conn.Close()
}
