package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jpfigueredo/rate-limiter-challenge/internal/adapter/http"
	"github.com/jpfigueredo/rate-limiter-challenge/internal/config"
	"github.com/jpfigueredo/rate-limiter-challenge/internal/usecase"
	"github.com/jpfigueredo/rate-limiter-challenge/pkg/storage"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	repo := storage.NewRedisRateLimiter(cfg.RedisAddr)
	if err := repo.Client.Ping(context.Background()).Err(); err != nil {
		panic("Falha ao conectar Redis: " + err.Error())
	}

	uc := usecase.NewRateLimiterUseCase(repo, cfg.MaxRequests, cfg.MaxTokenRequests, time.Second, cfg.BlockDuration)

	r := gin.Default()
	r.Use(http.RateLimiterMiddleware(uc))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.Run(":8080")
}
