package internal

import (
	"net/http"
	"sync"
	"time"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
}

// Worker executa as requisições concorrentes
func Worker(url string, requests int, results chan<- Result, wg *sync.WaitGroup) {
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

// CalculateReport gera as métricas a partir dos resultados
// func CalculateReport(results []Result, totalDuration time.Duration, detailed bool) Report {
// 	statusCodes := make(map[int]int)
// 	var latencies []time.Duration
// 	var successCount int

// 	for _, r := range results {
// 		statusCodes[r.StatusCode]++
// 		if r.StatusCode == 200 {
// 			successCount++
// 		}
// 		latencies = append(latencies, r.Duration)
// 	}

// 	sort.Slice(latencies, func(i, j int) bool {
// 		return latencies[i] < latencies[j]
// 	})

// 	min := latencies[0]
// 	max := latencies[len(latencies)-1]
// 	mean := time.Duration(0)
// 	for _, l := range latencies {
// 		mean += l
// 	}
// 	mean /= time.Duration(len(latencies))

// 	p95 := latencies[int(float64(len(latencies))*0.95)-1]
// 	p99 := latencies[int(float64(len(latencies))*0.99)-1]

// 	report := Report{
// 		TotalRequests: len(results),
// 		SuccessCount:  successCount,
// 		StatusCodes:   statusCodes,
// 		Duration:      totalDuration,
// 		Min:           min,
// 		Max:           max,
// 		Mean:          mean,
// 		P95:           p95,
// 		P99:           p99,
// 	}

// 	if detailed {
// 		report.Latencies = latencies
// 	}

// 	return report
// }
