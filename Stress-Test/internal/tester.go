package internal

import (
	"net/http"
	"sort"
	"sync"
	"time"
)

type Report struct {
	TotalRequests int
	SuccessCount  int
	StatusCodes   map[int]int
	Duration      time.Duration

	// Latency metrics
	Latencies []time.Duration
	Min       time.Duration
	Max       time.Duration
	Mean      time.Duration
	P95       time.Duration
	P99       time.Duration
}

func RunLoadTest(url string, totalRequests, concurrency int) Report {
	start := time.Now()

	jobs := make(chan int, totalRequests)
	responses := make(chan response, totalRequests)
	var wg sync.WaitGroup

	// Workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{
				Timeout: 5 * time.Second,
			}

			for range jobs {
				reqStart := time.Now()
				resp, err := client.Get(url)
				elapsed := time.Since(reqStart)

				if err != nil {
					responses <- response{status: -1, latency: elapsed}
					continue
				}
				responses <- response{status: resp.StatusCode, latency: elapsed}
				resp.Body.Close()
			}
		}()
	}

	// Distribui jobs
	go func() {
		for i := 0; i < totalRequests; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Aguarda workers e fecha responses
	go func() {
		wg.Wait()
		close(responses)
	}()

	// Coleta resultados
	statusCount := make(map[int]int)
	success := 0
	var latencies []time.Duration

	for res := range responses {
		if res.status == 200 {
			success++
		}
		statusCount[res.status]++
		latencies = append(latencies, res.latency)
	}

	// Estatísticas de latência
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	min := latencies[0]
	max := latencies[len(latencies)-1]
	mean := calcMean(latencies)
	p95 := percentile(latencies, 95)
	p99 := percentile(latencies, 99)

	return Report{
		TotalRequests: totalRequests,
		SuccessCount:  success,
		StatusCodes:   statusCount,
		Duration:      time.Since(start),
		Latencies:     latencies,
		Min:           min,
		Max:           max,
		Mean:          mean,
		P95:           p95,
		P99:           p99,
	}
}

type response struct {
	status  int
	latency time.Duration
}

func calcMean(latencies []time.Duration) time.Duration {
	var total time.Duration
	for _, l := range latencies {
		total += l
	}
	return total / time.Duration(len(latencies))
}

func percentile(latencies []time.Duration, p int) time.Duration {
	if len(latencies) == 0 {
		return 0
	}
	index := (p * len(latencies)) / 100
	if index >= len(latencies) {
		index = len(latencies) - 1
	}
	return latencies[index]
}
