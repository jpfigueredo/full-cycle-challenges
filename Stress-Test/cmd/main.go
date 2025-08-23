package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"net/http"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
}

type Report struct {
	TotalRequests int             `json:"TotalRequests"`
	SuccessCount  int             `json:"SuccessCount"`
	StatusCodes   map[int]int     `json:"StatusCodes"`
	Duration      time.Duration   `json:"Duration"`
	Latencies     []time.Duration `json:"Latencies,omitempty"`
	Min           time.Duration   `json:"Min"`
	Max           time.Duration   `json:"Max"`
	Mean          time.Duration   `json:"Mean"`
	P95           time.Duration   `json:"P95"`
	P99           time.Duration   `json:"P99"`
}

// Bucket representa um intervalo de latência
type Bucket struct {
	From  time.Duration `json:"From"`
	To    time.Duration `json:"To"`
	Count int           `json:"Count"`
}

func worker(url string, requests int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}
	for i := 0; i < requests; i++ {
		start := time.Now()
		resp, err := client.Get(url)
		duration := time.Since(start)
		if err != nil {
			results <- Result{StatusCode: -1, Duration: duration}
			continue
		}
		results <- Result{StatusCode: resp.StatusCode, Duration: duration}
		resp.Body.Close()
	}
}

func calculateReport(results []Result, totalDuration time.Duration, detailed bool) Report {
	statusCodes := make(map[int]int)
	var latencies []time.Duration
	var successCount int

	for _, r := range results {
		statusCodes[r.StatusCode]++
		if r.StatusCode == 200 {
			successCount++
		}
		latencies = append(latencies, r.Duration)
	}

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	min := latencies[0]
	max := latencies[len(latencies)-1]
	mean := time.Duration(0)
	for _, l := range latencies {
		mean += l
	}
	mean /= time.Duration(len(latencies))

	p95 := latencies[int(float64(len(latencies))*0.95)-1]
	p99 := latencies[int(float64(len(latencies))*0.99)-1]

	report := Report{
		TotalRequests: len(results),
		SuccessCount:  successCount,
		StatusCodes:   statusCodes,
		Duration:      totalDuration,
		Min:           min,
		Max:           max,
		Mean:          mean,
		P95:           p95,
		P99:           p99,
	}

	// só inclui as latências se --detailed for true
	if detailed {
		report.Latencies = latencies
	}

	return report
}

func printReport(report Report, output string) {
	switch strings.ToLower(output) {
	case "json":
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))
	case "csv":
		writer := csv.NewWriter(os.Stdout)
		defer writer.Flush()
		writer.Write([]string{"TotalRequests", "Successful", "Duration", "Min", "Max", "Mean", "P95", "P99", "StatusCodes"})
		var statusCodes []string
		for code, count := range report.StatusCodes {
			statusCodes = append(statusCodes, fmt.Sprintf("%d:%d", code, count))
		}
		writer.Write([]string{
			strconv.Itoa(report.TotalRequests),
			strconv.Itoa(report.SuccessCount),
			report.Duration.String(),
			report.Min.String(),
			report.Max.String(),
			report.Mean.String(),
			report.P95.String(),
			report.P99.String(),
			strings.Join(statusCodes, " "),
		})
	default:
		fmt.Println("===== Stress Test Report =====")
		fmt.Printf("Total Requests: %d\n", report.TotalRequests)
		fmt.Printf("Successful (200): %d\n", report.SuccessCount)
		fmt.Printf("Duration: %v\n\n", report.Duration)

		fmt.Println("Latency Metrics:")
		fmt.Printf("  Min:   %v\n", report.Min)
		fmt.Printf("  Max:   %v\n", report.Max)
		fmt.Printf("  Mean:  %v\n", report.Mean)
		fmt.Printf("  P95:   %v\n", report.P95)
		fmt.Printf("  P99:   %v\n", report.P99)
		fmt.Println()
		fmt.Println("Status Codes Distribution:")
		for code, count := range report.StatusCodes {
			fmt.Printf("  %d : %d\n", code, count)
		}
	}
}

// BuildHistogram constrói histograma de latências
func BuildHistogram(latencies []time.Duration, bucketSize time.Duration) ([]Bucket, string) {
	if len(latencies) == 0 {
		return nil, "No latencies recorded"
	}

	// encontrar a maior latência
	var maxLatency time.Duration
	for _, l := range latencies {
		if l > maxLatency {
			maxLatency = l
		}
	}

	// número de buckets
	numBuckets := int(maxLatency/bucketSize) + 1
	buckets := make([]Bucket, numBuckets)

	// inicializar buckets
	for i := 0; i < numBuckets; i++ {
		from := time.Duration(i) * bucketSize
		to := from + bucketSize
		buckets[i] = Bucket{From: from, To: to, Count: 0}
	}

	// distribuir latências
	for _, l := range latencies {
		index := int(l / bucketSize)
		if index >= len(buckets) {
			index = len(buckets) - 1
		}
		buckets[index].Count++
	}

	// gerar string legível
	var sb strings.Builder
	sb.WriteString("Latency Histogram:\n")
	for _, b := range buckets {
		sb.WriteString(fmt.Sprintf("  [%v - %v] : %d\n", b.From, b.To, b.Count))
	}

	return buckets, sb.String()
}

func main() {
	url := flag.String("url", "", "Target URL to test")
	requests := flag.Int("requests", 1, "Total number of requests to perform")
	concurrency := flag.Int("concurrency", 1, "Number of concurrent workers")
	output := flag.String("output", "default", "Output format: default, json, csv")
	detailed := flag.Bool("detailed", false, "Include detailed latencies in JSON output")
	flag.Parse()

	if *url == "" {
		fmt.Println("Error: --url is required")
		os.Exit(1)
	}

	results := make(chan Result, *requests)
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		reqs := *requests / *concurrency
		if i < *requests%*concurrency {
			reqs++
		}
		go worker(*url, reqs, results, &wg)
	}

	wg.Wait()
	close(results)

	var allResults []Result
	for r := range results {
		allResults = append(allResults, r)
	}

	totalDuration := time.Since(start)
	report := calculateReport(allResults, totalDuration, *detailed)
	printReport(report, *output)

	latencies := []time.Duration{
		50 * time.Millisecond,
		120 * time.Millisecond,
		180 * time.Millisecond,
		220 * time.Millisecond,
		900 * time.Millisecond,
	}
	_, text := BuildHistogram(latencies, 100*time.Millisecond)
	fmt.Println(text)

}
