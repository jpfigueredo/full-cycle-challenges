package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jpfigueredo/cep-clima-distributed/internal/shared/entity"
	"github.com/jpfigueredo/cep-clima-distributed/internal/shared/usecase"
)

type mockCEPRepo struct{}

func (m *mockCEPRepo) GetLocalizacao(ctx context.Context, cep string) (*entity.Localizacao, error) {
	if cep == "12345678" {
		return &entity.Localizacao{Localidade: "Test City"}, nil
	}
	return &entity.Localizacao{}, nil
}

type mockClimaRepo struct{}

func (m *mockClimaRepo) GetClimaAtual(ctx context.Context, localidade string) (*entity.ClimaInput, error) {
	return &entity.ClimaInput{Temp: 25.5}, nil
}

func TestValidacaoCEPInvalido(t *testing.T) {
	uc := usecase.NewCEPClimaUseCase(&mockCEPRepo{}, &mockClimaRepo{})
	_, err := uc.CheckCEPAndGetClima(context.Background(), "1234-567") // Len <8 após clean
	if err == nil || err != usecase.ErrInvalidCEP {
		t.Error("Expected ErrInvalidCEP")
	}
}

func TestCEPNotFound(t *testing.T) {
	uc := usecase.NewCEPClimaUseCase(&mockCEPRepo{}, &mockClimaRepo{})
	_, err := uc.CheckCEPAndGetClima(context.Background(), "00000000")
	if err == nil || err != usecase.ErrCEPNotFound {
		t.Error("Expected ErrCEPNotFound")
	}
}

func TestSucessoComConversoes(t *testing.T) {
	uc := usecase.NewCEPClimaUseCase(&mockCEPRepo{}, &mockClimaRepo{})
	out, err := uc.CheckCEPAndGetClima(context.Background(), "12345678")
	if err != nil {
		t.Error("Expected no error")
	}
	if out.TempC != 25.5 || out.TempF != 77.9 || out.TempK != 298.5 {
		t.Error("Expected conversões corretas")
	}
}

func TestClimaErro(t *testing.T) {
	mockClimaErr := &mockClimaRepoErr{}
	uc := usecase.NewCEPClimaUseCase(&mockCEPRepo{}, mockClimaErr)
	_, err := uc.CheckCEPAndGetClima(context.Background(), "12345678")
	if err == nil || err.Error() != "clima fail" {
		t.Error("Expected propagated clima error")
	}
}

type mockClimaRepoErr struct{}

func (m *mockClimaRepoErr) GetClimaAtual(ctx context.Context, localidade string) (*entity.ClimaInput, error) {
	return nil, errors.New("clima fail")
}
