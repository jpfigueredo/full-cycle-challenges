package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jpfigueredo/rate-limiter-challenge/internal/entity"
	"github.com/jpfigueredo/rate-limiter-challenge/internal/repository"
	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	Client *redis.Client
}

var _ repository.RateLimiterRepository = (*RedisRateLimiter)(nil)

func NewRedisRateLimiter(addr string) *RedisRateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisRateLimiter{Client: client}
}

func (r *RedisRateLimiter) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	rateKey := fmt.Sprintf("rate:%s", key)
	val, err := r.Client.Incr(ctx, rateKey).Result()
	if err != nil {
		return 0, err
	}
	if val == 1 {
		r.Client.Expire(ctx, rateKey, window)
	}
	return val, nil
}

func (r *RedisRateLimiter) Block(ctx context.Context, key string, blockDuration time.Duration) error {
	blockKey := fmt.Sprintf("block:%s", key)
	return r.Client.Set(ctx, blockKey, "blocked", blockDuration).Err()
}

func (r *RedisRateLimiter) IsBlocked(ctx context.Context, key string) (bool, error) {
	blockKey := fmt.Sprintf("block:%s", key)
	_, err := r.Client.Get(ctx, blockKey).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisRateLimiter) GetState(ctx context.Context, key string) (*entity.RateLimit, error) {
	rateKey := fmt.Sprintf("rate:%s", key)
	blockKey := fmt.Sprintf("block:%s", key)

	count, err := r.Client.Get(ctx, rateKey).Int64()
	if err == redis.Nil {
		count = 0
	} else if err != nil {
		return nil, err
	}

	ttl, err := r.Client.TTL(ctx, blockKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	blockedUntil := time.Time{}
	if ttl > 0 {
		blockedUntil = time.Now().Add(ttl)
	}

	return &entity.RateLimit{
		Key:          key,
		Count:        count,
		BlockedUntil: blockedUntil,
	}, nil
}
