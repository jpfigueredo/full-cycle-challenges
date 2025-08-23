package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/jpfigueredo/rate-limiter-challenge/internal/config"
)

func TestLoadAllEnvs(t *testing.T) {
	os.Setenv("MAX_REQUESTS_PER_SECOND", "3")
	os.Setenv("MAX_TOKEN_REQUESTS_PER_SECOND", "7")
	os.Setenv("WINDOW_SECONDS", "3")
	os.Setenv("BLOCK_DURATION_SECONDS", "600")
	os.Setenv("REDIS_ADDR", "test:6379")
	defer os.Clearenv()

	cfg := config.Load()
	if cfg.MaxRequests != 3 || cfg.MaxTokenRequests != 7 || cfg.Window != 3*time.Second || cfg.BlockDuration != 600*time.Second || cfg.RedisAddr != "test:6379" {
		t.Error("Expected all envs loaded")
	}
}

func TestLoadPartialEnvs(t *testing.T) {
	os.Setenv("MAX_REQUESTS_PER_SECOND", "4")
	defer os.Clearenv()

	cfg := config.Load()
	if cfg.MaxRequests != 4 || cfg.Window != time.Second || cfg.RedisAddr != "localhost:6379" {
		t.Error("Expected partial envs with defaults")
	}
}
