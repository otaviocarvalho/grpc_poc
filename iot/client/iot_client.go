package main

import (
    "log"
    "time"
    "flag"
    "math/rand"
    "sync"

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
    wg.Add(n * m)

    var histMutex sync.RWMutex

    startTime := time.Now()
    for i := 0; i < n; i++ {
        go func() {

            for j := 0; j < m; j++ {
                startTimeLoop := time.Now()

                // Controls processing local and remote
                if ((j+1) % *batchSize != 0) {
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

    plotQuantiles(n*m, totalTime, hist)

    conn.Close()
}
