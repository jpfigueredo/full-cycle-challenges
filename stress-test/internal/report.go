package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Bucket struct {
	From  time.Duration `json:"From"`
	To    time.Duration `json:"To"`
	Count int           `json:"Count"`
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
	Histogram     []Bucket        `json:"Histogram,omitempty"`
}

func CalculateReport(results []Result, totalDuration time.Duration, detailed bool) Report {
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

	if detailed {
		report.Latencies = latencies
	}

	return report
}

func PrintReport(report Report, output string) {
	switch strings.ToLower(output) {
	case "json":
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))

	case "csv":
		writer := csv.NewWriter(os.Stdout)
		defer writer.Flush()
		writer.Write([]string{
			"TotalRequests", "Successful", "Duration", "Min", "Max", "Mean", "P95", "P99", "StatusCodes",
		})
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

		if len(report.Histogram) > 0 {
			fmt.Println()
			fmt.Println("Latency Histogram:")
			for _, b := range report.Histogram {
				fmt.Printf("  [%v - %v] : %d\n", b.From, b.To, b.Count)
			}
		}
	}
}
