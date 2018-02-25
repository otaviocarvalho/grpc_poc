package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var inputDir = flag.String("i", "./", "input dir")
var outputFile = flag.String("o", "./output.csv", "output file")

type Stats struct {
	Qps       float64 `json:"qps"`
	Perc50    int64   `json:"p50"`
	Perc90    int64   `json:"p90"`
	Perc99    int64   `json:"p99"`
	Perc999   int64   `json:"p999"`
	Perc9999  int64   `json:"p9999"`
	Perc99999 int64   `json:"p99999"`
}

func main() {
	flag.Parse()

	files, err := ioutil.ReadDir(*inputDir)
	if err != nil {
		fmt.Println("Reading dir", err.Error())
		os.Exit(1)
	}

	// Read from json input
	var stats Stats
	var listStats []Stats
	for _, f := range files {
		raw, err := ioutil.ReadFile(fmt.Sprintf("%v%v", *inputDir, f.Name()))
		if err != nil {
			fmt.Println("Reading from input file", err.Error())
			os.Exit(1)
		}

		json.Unmarshal(raw, &stats)
		listStats = append(listStats, stats)
	}

	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)

	// Write header
	var header = []string{
		"qps",
		"p50",
		"p90",
		"p99",
		"p999",
		"p9999",
		"p99999",
	}
	if err = writer.Write(header); err != nil {
		fmt.Println("Error writing to output file", err.Error())
	}

	// Write body
	for _, entryStats := range listStats {
		var line = []string{
			fmt.Sprint(entryStats.Qps),
			fmt.Sprint(entryStats.Perc50),
			fmt.Sprint(entryStats.Perc90),
			fmt.Sprint(entryStats.Perc99),
			fmt.Sprint(entryStats.Perc999),
			fmt.Sprint(entryStats.Perc9999),
			fmt.Sprint(entryStats.Perc99999),
		}

		if err = writer.Write(line); err != nil {
			fmt.Println("Error writing to output buffer", err.Error())
		}
	}
	writer.Flush()

	// Write to csv output
	if err = ioutil.WriteFile(*outputFile, buffer.Bytes(), 0644); err != nil {
		fmt.Println("Writing to output file", err.Error())
		os.Exit(1)
	}
}
