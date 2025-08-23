package entity

import "time"

type RateLimit struct {
	Key          string
	Count        int64
	BlockedUntil time.Time
}
