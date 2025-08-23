package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	middleware "github.com/jpfigueredo/rate-limiter-challenge/internal/adapter/http"
	"github.com/jpfigueredo/rate-limiter-challenge/internal/entity"
	"github.com/jpfigueredo/rate-limiter-challenge/internal/usecase"
)

type mockRepo struct{} // Mock local para testes independentes

func (m *mockRepo) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	return 1, nil
}
func (m *mockRepo) Block(ctx context.Context, key string, blockDuration time.Duration) error {
	return nil
}
func (m *mockRepo) IsBlocked(ctx context.Context, key string) (bool, error) { return false, nil }
func (m *mockRepo) GetState(ctx context.Context, key string) (*entity.RateLimit, error) {
	return &entity.RateLimit{}, nil
}

func TestMiddlewareAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	uc := usecase.NewRateLimiterUseCase(&mockRepo{}, 5, 10, time.Second, 5*time.Minute)
	r.Use(middleware.RateLimiterMiddleware(uc))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Error("Expected allowed")
	}
}

func TestMiddlewareBlocked(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mock := &mockRepoBlocked{}
	uc := usecase.NewRateLimiterUseCase(mock, 1, 10, time.Second, 5*time.Minute)
	r.Use(middleware.RateLimiterMiddleware(uc))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 429 {
		t.Error("Expected blocked with 429")
	}
}

type mockRepoBlocked struct{}

func (m *mockRepoBlocked) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	return 2, nil
} // > max
func (m *mockRepoBlocked) Block(ctx context.Context, key string, blockDuration time.Duration) error {
	return nil
}
func (m *mockRepoBlocked) IsBlocked(ctx context.Context, key string) (bool, error) { return false, nil }
func (m *mockRepoBlocked) GetState(ctx context.Context, key string) (*entity.RateLimit, error) {
	return &entity.RateLimit{}, nil
}
