package repository

import (
	"context"
	"time"

	"github.com/jpfigueredo/rate-limiter-challenge/internal/entity"
)

type RateLimiterRepository interface {
	Increment(ctx context.Context, key string, window time.Duration) (int64, error)
	Block(ctx context.Context, key string, blockDuration time.Duration) error
	IsBlocked(ctx context.Context, key string) (bool, error)
	GetState(ctx context.Context, key string) (*entity.RateLimit, error)
}
