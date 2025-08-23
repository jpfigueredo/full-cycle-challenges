package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/jpfigueredo/rate-limiter-challenge/pkg/storage"
	"github.com/redis/go-redis/v9"
)

func TestIncrementMultiple(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	repo := &storage.RedisRateLimiter{Client: client}

	count, err := repo.Increment(context.Background(), "test", time.Second)
	if err != nil || count != 1 {
		t.Error("Expected 1 on first")
	}

	count, err = repo.Increment(context.Background(), "test", time.Second)
	if err != nil || count != 2 {
		t.Error("Expected 2 on second")
	}

	mr.FastForward(time.Second + time.Millisecond)
	count, err = repo.Increment(context.Background(), "test", time.Second)
	if err != nil || count != 1 {
		t.Error("Expected reset to 1 after expiration")
	}
}

func TestBlockAndIsBlocked(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	repo := &storage.RedisRateLimiter{Client: client}

	err := repo.Block(context.Background(), "test", 10*time.Second)
	if err != nil {
		t.Error("Expected no error on block")
	}

	blocked, err := repo.IsBlocked(context.Background(), "test")
	if err != nil || !blocked {
		t.Error("Expected blocked")
	}

	mr.FastForward(11 * time.Second)
	blocked, err = repo.IsBlocked(context.Background(), "test")
	if err != nil || blocked {
		t.Error("Expected not blocked after expiration")
	}
}

func TestGetStateNoKey(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	repo := &storage.RedisRateLimiter{Client: client}

	state, err := repo.GetState(context.Background(), "nonexistent")
	if err != nil || state.Count != 0 || !state.BlockedUntil.IsZero() {
		t.Error("Expected zero state for nonexistent key")
	}
}

func TestIncrementError(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	repo := &storage.RedisRateLimiter{Client: client}
	client.Close()

	_, err := repo.Increment(context.Background(), "test", time.Second)
	if err == nil {
		t.Error("Expected error on increment with closed client")
	}
}

func TestBlockError(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	repo := &storage.RedisRateLimiter{Client: client}
	client.Close()

	err := repo.Block(context.Background(), "test", time.Second)
	if err == nil {
		t.Error("Expected error on block with closed client")
	}
}

func TestIsBlockedError(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	repo := &storage.RedisRateLimiter{Client: client}
	client.Close()

	_, err := repo.IsBlocked(context.Background(), "test")
	if err == nil {
		t.Error("Expected error on isblocked with closed client")
	}
}

func TestIncrementPartialExpiration(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	repo := &storage.RedisRateLimiter{Client: client}

	_, _ = repo.Increment(context.Background(), "test", 10*time.Second)
	mr.FastForward(5 * time.Second)
	count, err := repo.Increment(context.Background(), "test", 10*time.Second)
	if err != nil || count != 2 {
		t.Error("Expected count 2 after partial time")
	}
}

// func TestBlockZeroDuration(t *testing.T) {
// 	mr := miniredis.RunT(t)
// 	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
// 	repo := &storage.RedisRateLimiter{Client: client}

// 	err := repo.Block(context.Background(), "test", 0)
// 	if err != nil {
// 		t.Errorf("Expected no error on zero duration block: got %v", err)
// 	}

// 	mr.FastForward(time.Second)
// 	blocked, err := repo.IsBlocked(context.Background(), "test")
// 	if err != nil || blocked {
// 		t.Errorf("Expected not blocked with zero duration: got blocked=%v, err=%v", blocked, err)
// 	}
// }

func TestGetStateFull(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	repo := &storage.RedisRateLimiter{Client: client}

	_, _ = repo.Increment(context.Background(), "test", 20*time.Second)
	_ = repo.Block(context.Background(), "test", 10*time.Second)

	state, err := repo.GetState(context.Background(), "test")
	if err != nil || state.Count != 1 || state.BlockedUntil.IsZero() {
		t.Errorf("Expected count=1 and BlockedUntil not zero: got Count=%d, BlockedUntil=%v", state.Count, state.BlockedUntil)
	}

	mr.FastForward(11 * time.Second)
	state, err = repo.GetState(context.Background(), "test")
	if err != nil || state.Count != 1 || !state.BlockedUntil.IsZero() {
		t.Errorf("Expected count=1 persist and BlockedUntil zero: got Count=%d, BlockedUntil=%v", state.Count, state.BlockedUntil)
	}
}

func TestGetStateError(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	repo := &storage.RedisRateLimiter{Client: client}
	client.Close()

	_, err := repo.GetState(context.Background(), "test")
	if err == nil {
		t.Error("Expected error on getstate with closed client")
	}
}
