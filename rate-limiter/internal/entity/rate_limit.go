package entity

import "time"

type RateLimit struct {
	Key          string // IP ou Token
	Count        int64
	BlockedUntil time.Time
}
