package internal

import (
	"fmt"
	"strings"
	"time"
)

func BuildHistogram(latencies []time.Duration, bucketSize time.Duration) ([]Bucket, string) {
	if len(latencies) == 0 {
		return nil, "No latencies recorded"
	}

	var maxLatency time.Duration
	for _, l := range latencies {
		if l > maxLatency {
			maxLatency = l
		}
	}

	numBuckets := int(maxLatency/bucketSize) + 1
	buckets := make([]Bucket, numBuckets)
	for i := 0; i < numBuckets; i++ {
		from := time.Duration(i) * bucketSize
		to := from + bucketSize
		buckets[i] = Bucket{From: from, To: to, Count: 0}
	}

	for _, l := range latencies {
		idx := int(l / bucketSize)
		if idx >= len(buckets) {
			idx = len(buckets) - 1
		}
		buckets[idx].Count++
	}

	var sb strings.Builder
	sb.WriteString("Latency Histogram:\n")
	for _, b := range buckets {
		sb.WriteString(fmt.Sprintf("  [%v - %v] : %d\n", b.From, b.To, b.Count))
	}

	return buckets, sb.String()
}
