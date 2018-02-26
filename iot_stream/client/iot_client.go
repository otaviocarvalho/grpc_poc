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

	pb "grpc_poc/iot_stream"

	hdr "github.com/otaviocarvalho/hdrhistogram"
)

var total = flag.Int("n", 1, "total requests for all clients")
var host = flag.String("s", "localhost:50051", "ip/port")
var outputFile = flag.String("o", "./output_stats.json", "output json file")

type Stats struct {
	Qps       float64 `json:"qps"`
	Perc50    int64   `json:"p50"`
	Perc90    int64   `json:"p90"`
	Perc99    int64   `json:"p99"`
	Perc999   int64   `json:"p999"`
	Perc9999  int64   `json:"p9999"`
	Perc99999 int64   `json:"p99999"`
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

	fmt.Printf("%+v", stats)
}

func main() {
	flag.Parse()

	// 1ms to 30 seconds range, 5 sigfigs precision
	hist := hdr.New(1000000, 30000000000, 5)

	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(*host,
		grpc.WithInsecure(),
		grpc.WithMaxMsgSize(64<<20),
	)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Creates a new data client
	client := pb.NewMeterCommunicatorClient(conn)
	stream, err := client.SimpleRPC(context.Background())

	data := &pb.Measurement{
		Value: rand.Int63(),
	}

	var histMutex sync.RWMutex

	// Sends measurement
	startTime := time.Now()
	for i := 0; i < *total; i++ {

		startTimeLoop := time.Now()

		err := stream.Send(data)
		if err != nil {
			log.Fatalf("error sending msg: %v", err)
		}

		totalTimeLoop := time.Now().Sub(startTimeLoop)

		histMutex.Lock()
		hist.RecordValue(totalTimeLoop.Nanoseconds())
		histMutex.Unlock()

	}

	totalTime := time.Now().Sub(startTime)

	plotStats(*total, totalTime, hist)
	saveStats(*total, totalTime, hist)

	stream.CloseSend()
}
