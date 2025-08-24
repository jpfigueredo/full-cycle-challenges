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
