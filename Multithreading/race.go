package main

import (
	"context"
	"fmt"
	"time"
)

type fetchResult struct {
	addr Address
	err  error
}

// RaceService corre as requisições nas APIs e retorna o resultado mais rápido válido.
type RaceService struct {
	fetchers []AddressFetcher
	timeout  time.Duration
}

func NewRaceService(fetchers []AddressFetcher, timeout time.Duration) *RaceService {
	return &RaceService{
		fetchers: fetchers,
		timeout:  timeout,
	}
}

// Run dispara as requisições simultâneas e retorna o primeiro sucesso ou erro.
func (r *RaceService) Run(ctx context.Context, cep string) (Address, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	results := make(chan fetchResult, len(r.fetchers))
	childCancels := make([]context.CancelFunc, len(r.fetchers))

	for i, f := range r.fetchers {
		childCtx, childCancel := context.WithCancel(ctx)
		childCancels[i] = childCancel

		go func(fetcher AddressFetcher, c context.Context) {
			addr, err := fetcher.Fetch(c, cep)
			results <- fetchResult{addr: addr, err: err}
		}(f, childCtx)
	}

	var firstErr error
	for i := 0; i < len(r.fetchers); i++ {
		select {
		case res := <-results:
			if res.err == nil {
				// sucesso, cancelar os outros
				for _, cancel := range childCancels {
					cancel()
				}
				return res.addr, nil
			}
			if firstErr == nil {
				firstErr = res.err
			}
		case <-ctx.Done():
			return Address{}, fmt.Errorf("timeout after %s", r.timeout)
		}
	}
	return Address{}, firstErr
}
