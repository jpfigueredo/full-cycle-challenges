package usecase

import (
	"context"
	"time"

	"github.com/jpfigueredo/rate-limiter-challenge/internal/entity"
	"github.com/jpfigueredo/rate-limiter-challenge/internal/repository"
)

type RateLimiterUseCase struct {
	Repo          repository.RateLimiterRepository
	MaxRequests   int64         // Ex.: 5 por IP, carregado de config
	MaxTokenReqs  int64         // Sobrepõe se token presente
	Window        time.Duration // Janela de tempo, ex.: 1s
	BlockDuration time.Duration // Tempo de bloqueio, ex.: 5min
}

func NewRateLimiterUseCase(repo repository.RateLimiterRepository, maxReq int64, maxToken int64, window, block time.Duration) *RateLimiterUseCase {
	return &RateLimiterUseCase{
		Repo:          repo,
		MaxRequests:   maxReq,
		MaxTokenReqs:  maxToken,
		Window:        window,
		BlockDuration: block,
	}
}

func (uc *RateLimiterUseCase) CheckAndIncrement(ctx context.Context, ip, token string) (bool, error) {
	key := ip
	max := uc.MaxRequests
	if token != "" {
		key = "token:" + token // Prefixo para diferenciar
		max = uc.MaxTokenReqs
	}

	// Verifica bloqueio primeiro
	blocked, err := uc.Repo.IsBlocked(ctx, key)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil // Bloqueado: nega acesso
	}

	// Incrementa e verifica limite
	count, err := uc.Repo.Increment(ctx, key, uc.Window)
	if err != nil {
		return false, err
	}
	if count > max {
		// Bloqueia se excedido
		if err := uc.Repo.Block(ctx, key, uc.BlockDuration); err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil // Permitido
}

func (uc *RateLimiterUseCase) GetLimitState(ctx context.Context, ip, token string) (*entity.RateLimit, error) {
	key := ip
	if token != "" {
		key = "token:" + token
	}
	return uc.Repo.GetState(ctx, key)
}
