package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jpfigueredo/cep-clima-distributed/service-a/internal/entity"
	"github.com/jpfigueredo/cep-clima-distributed/service-a/internal/usecase"
)

type mockCEPRepoError struct{}
type mockClimaRepo struct{}

func (m *mockCEPRepoError) GetLocalizacao(ctx context.Context, cep string) (*entity.Localizacao, error) {
	return nil, errors.New("repo fail")
}

func (m *mockClimaRepo) GetClimaAtual(ctx context.Context, localidade string) (*entity.ClimaInput, error) {
	return &entity.ClimaInput{Temp: 25.0}, nil
}

func TestRepoError(t *testing.T) {
	mockCEP := &mockCEPRepoError{}
	uc := usecase.NewCEPClimaUseCase(mockCEP, &mockClimaRepo{})
	_, err := uc.CheckCEPAndGetClima(context.Background(), "12345678")
	if err == nil || err.Error() != "repo fail" {
		t.Error("Expected propagated repo error")
	}
}
