package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jpfigueredo/rate-limiter-challenge/internal/entity"
	"github.com/jpfigueredo/rate-limiter-challenge/internal/usecase"
)

type mockRepo struct {
	count         int64
	blocked       bool
	blockErr      error
	incErr        error
	blockCheckErr error
	state         *entity.RateLimit
	stateErr      error
}

func (m *mockRepo) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	return m.count, m.incErr
}
func (m *mockRepo) Block(ctx context.Context, key string, blockDuration time.Duration) error {
	return m.blockErr
}
func (m *mockRepo) IsBlocked(ctx context.Context, key string) (bool, error) {
	return m.blocked, m.blockCheckErr
}
func (m *mockRepo) GetState(ctx context.Context, key string) (*entity.RateLimit, error) {
	return &entity.RateLimit{Key: key, Count: m.count, BlockedUntil: time.Time{}}, m.stateErr
}

func TestCheckAndIncrementIPUnderLimit(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{count: 4}, 5, 10, time.Second, 5*time.Minute)
	allowed, err := uc.CheckAndIncrement(context.Background(), "127.0.0.1", "")
	if err != nil || !allowed {
		t.Error("Expected allowed for IP under limit")
	}
}

func TestCheckAndIncrementIPExceed(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{count: 6}, 5, 10, time.Second, 5*time.Minute)
	allowed, err := uc.CheckAndIncrement(context.Background(), "127.0.0.1", "")
	if err != nil || allowed {
		t.Error("Expected denied for IP exceed")
	}
}

func TestCheckAndIncrementTokenPriority(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{count: 6}, 5, 10, time.Second, 5*time.Minute)
	allowed, err := uc.CheckAndIncrement(context.Background(), "127.0.0.1", "mytoken")
	if err != nil || !allowed {
		t.Error("Expected allowed for token over IP limit")
	}
}

func TestCheckAndIncrementTokenExceed(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{count: 11}, 5, 10, time.Second, 5*time.Minute)
	allowed, err := uc.CheckAndIncrement(context.Background(), "127.0.0.1", "mytoken")
	if err != nil || allowed {
		t.Error("Expected denied for token exceed")
	}
}

func TestCheckAndIncrementBlocked(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{blocked: true}, 5, 10, time.Second, 5*time.Minute)
	allowed, err := uc.CheckAndIncrement(context.Background(), "127.0.0.1", "")
	if err != nil || allowed {
		t.Error("Expected denied when already blocked")
	}
}

func TestCheckAndIncrementBlockError(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{count: 6, blockErr: errors.New("block fail")}, 5, 10, time.Second, 5*time.Minute)
	_, err := uc.CheckAndIncrement(context.Background(), "127.0.0.1", "")
	if err == nil {
		t.Error("Expected error when block fails")
	}
}

func TestCheckAndIncrementIncError(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{incErr: errors.New("inc fail")}, 5, 10, time.Second, 5*time.Minute)
	_, err := uc.CheckAndIncrement(context.Background(), "127.0.0.1", "")
	if err == nil {
		t.Error("Expected error when increment fails")
	}
}

func TestGetLimitState(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{count: 3}, 5, 10, time.Second, 5*time.Minute)
	state, err := uc.GetLimitState(context.Background(), "127.0.0.1", "")
	if err != nil || state.Key != "127.0.0.1" || state.Count != 3 {
		t.Errorf("Expected correct state for IP: got Key=%s, Count=%d, err=%v", state.Key, state.Count, err)
	}

	state, err = uc.GetLimitState(context.Background(), "127.0.0.1", "mytoken")
	if err != nil || state.Key != "token:mytoken" || state.Count != 3 {
		t.Errorf("Expected token key priority: got Key=%s, Count=%d, err=%v", state.Key, state.Count, err)
	}
}
func TestGetLimitStateError(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{stateErr: errors.New("state fail")}, 5, 10, time.Second, 5*time.Minute)
	_, err := uc.GetLimitState(context.Background(), "127.0.0.1", "")
	if err == nil {
		t.Error("Expected error when get state fails")
	}
}

func TestCheckAndIncrementAtLimit(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{count: 5}, 5, 10, time.Second, 5*time.Minute)
	allowed, err := uc.CheckAndIncrement(context.Background(), "127.0.0.1", "")
	if err != nil || !allowed {
		t.Error("Expected allowed at exact limit")
	}
}

func TestCheckAndIncrementIsBlockedError(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{blockCheckErr: errors.New("check fail")}, 5, 10, time.Second, 5*time.Minute)
	_, err := uc.CheckAndIncrement(context.Background(), "127.0.0.1", "")
	if err == nil {
		t.Error("Expected error when isblocked fails")
	}
}

func TestGetLimitStateToken(t *testing.T) {
	expected := &entity.RateLimit{Key: "token:mytoken", Count: 5}
	uc := usecase.NewRateLimiterUseCase(&mockRepo{state: expected}, 5, 10, time.Second, 5*time.Minute)
	state, err := uc.GetLimitState(context.Background(), "127.0.0.1", "mytoken")
	if err != nil || state.Key != "token:mytoken" {
		t.Error("Expected token key in state")
	}
}

func TestCheckAndIncrementZeroCount(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{count: 0}, 5, 10, time.Second, 5*time.Minute)
	allowed, err := uc.CheckAndIncrement(context.Background(), "127.0.0.1", "")
	if err != nil || !allowed {
		t.Error("Expected allowed with zero count")
	}
}

func TestCheckAndIncrementBlockedWithToken(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{blocked: true}, 5, 10, time.Second, 5*time.Minute)
	allowed, err := uc.CheckAndIncrement(context.Background(), "127.0.0.1", "mytoken")
	if err != nil || allowed {
		t.Error("Expected denied when blocked with token")
	}
}

func TestGetLimitStateZero(t *testing.T) {
	uc := usecase.NewRateLimiterUseCase(&mockRepo{state: &entity.RateLimit{Key: "127.0.0.1", Count: 0}}, 5, 10, time.Second, 5*time.Minute)
	state, err := uc.GetLimitState(context.Background(), "127.0.0.1", "")
	if err != nil || state.Count != 0 {
		t.Error("Expected zero count state")
	}
}
