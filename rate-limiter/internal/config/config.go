package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	MaxRequests      int64
	MaxTokenRequests int64
	Window           time.Duration
	BlockDuration    time.Duration
	RedisAddr        string
}

func Load() *Config {
	maxReq, _ := strconv.ParseInt(os.Getenv("MAX_REQUESTS_PER_SECOND"), 10, 64)
	maxToken, _ := strconv.ParseInt(os.Getenv("MAX_TOKEN_REQUESTS_PER_SECOND"), 10, 64)
	windowSec, _ := strconv.ParseInt(os.Getenv("WINDOW_SECONDS"), 10, 64)
	if windowSec == 0 {
		windowSec = 1
	}
	blockSec, _ := strconv.ParseInt(os.Getenv("BLOCK_DURATION_SECONDS"), 10, 64)
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	return &Config{
		MaxRequests:      maxReq,
		MaxTokenRequests: maxToken,
		Window:           time.Duration(windowSec) * time.Second,
		BlockDuration:    time.Duration(blockSec) * time.Second,
		RedisAddr:        redisAddr,
	}
}
