package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/jpfigueredo/cep-clima-distributed/internal/shared/entity"
	"github.com/jpfigueredo/cep-clima-distributed/internal/shared/repository"
)

var (
	ErrInvalidCEP    = errors.New("invalid zipcode")
	ErrCEPNotFound   = errors.New("can not find zipcode")
	ErrClimaNotFound = errors.New("clima not found")
)

type CEPClimaUseCase struct {
	CEPRepo   repository.CEPRepository
	ClimaRepo repository.ClimaRepository
}

func NewCEPClimaUseCase(cepRepo repository.CEPRepository, climaRepo repository.ClimaRepository) *CEPClimaUseCase {
	return &CEPClimaUseCase{CEPRepo: cepRepo, ClimaRepo: climaRepo}
}

func (uc *CEPClimaUseCase) CheckCEPAndGetClima(ctx context.Context, cep string) (*entity.ClimaOutput, error) {
	cleanCEP := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, cep)
	fmt.Printf("Debug: Input CEP=%s, Cleaned=%s, Len=%d\n", cep, cleanCEP, len(cleanCEP))
	if len(cleanCEP) != 8 {
		return nil, ErrInvalidCEP
	}

	loc, err := uc.CEPRepo.GetLocalizacao(ctx, cleanCEP)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Debug: Loc from ViaCEP=%+v, Err=%v\n", loc, err)
	if loc.Localidade == "" {
		return nil, ErrCEPNotFound
	}

	clima, err := uc.ClimaRepo.GetClimaAtual(ctx, loc.Localidade)
	fmt.Printf("Debug: Clima from WeatherAPI=%+v, Err=%v\n", clima, err)
	if err != nil {
		return nil, err
	}

	tempC := math.Round(clima.Temp*10) / 10
	tempF := math.Round((clima.Temp*1.8+32)*10) / 10
	tempK := math.Round((clima.Temp+273)*10) / 10

	return &entity.ClimaOutput{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}
