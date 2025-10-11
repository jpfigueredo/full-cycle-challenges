package repository

import (
	"context"

	"github.com/jpfigueredo/cep-clima-distributed/internal/shared/entity"
)

type CEPRepository interface {
	GetLocalizacao(ctx context.Context, cep string) (*entity.Localizacao, error)
}

type ClimaRepository interface {
	GetClimaAtual(ctx context.Context, localidade string) (*entity.ClimaInput, error)
}
