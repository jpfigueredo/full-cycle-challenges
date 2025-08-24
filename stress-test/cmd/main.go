package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"stress-test/internal"
)

func main() {
	url := flag.String("url", "", "Target URL to test")
	requests := flag.Int("requests", 1, "Total number of requests to perform")
	concurrency := flag.Int("concurrency", 1, "Number of concurrent workers")
	output := flag.String("output", "default", "Output format: default, json, csv")
	detailed := flag.Bool("detailed", false, "Include detailed latencies and histogram in output")
	flag.Parse()

	if *url == "" {
		fmt.Println("Error: --url is required")
		os.Exit(1)
	}

	results := make(chan internal.Result, *requests)
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		reqs := *requests / *concurrency
		if i < *requests%*concurrency {
			reqs++
		}
		go internal.Worker(*url, reqs, results, &wg)
	}

	wg.Wait()
	close(results)

	var allResults []internal.Result
	for r := range results {
		allResults = append(allResults, r)
	}

	totalDuration := time.Since(start)
	report := internal.CalculateReport(allResults, totalDuration, *detailed)

	if *detailed && len(report.Latencies) > 0 {
		hist, _ := internal.BuildHistogram(report.Latencies, 100*time.Millisecond)
		report.Histogram = hist
	}

	internal.PrintReport(report, *output)
}
