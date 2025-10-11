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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	ctx, span := otel.Tracer("service-b").Start(ctx, "CheckCEPAndGetClima")
	defer span.End()

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

	ctxVia, spanVia := otel.Tracer("service-b").Start(ctx, "GetLocalizacao")
	loc, err := uc.CEPRepo.GetLocalizacao(ctxVia, cleanCEP)
	spanVia.SetAttributes(attribute.String("cep", cleanCEP))
	spanVia.End()
	if err != nil {
		return nil, err
	}
	if loc.Localidade == "" {
		return nil, ErrCEPNotFound
	}

	ctxClima, spanClima := otel.Tracer("service-b").Start(ctx, "GetClimaAtual")
	clima, err := uc.ClimaRepo.GetClimaAtual(ctxClima, loc.Localidade)
	spanClima.SetAttributes(attribute.String("localidade", loc.Localidade))
	spanClima.End()
	if err != nil {
		return nil, err
	}

	tempC := math.Round(clima.Temp*10) / 10
	tempF := math.Round((clima.Temp*1.8+32)*10) / 10
	tempK := math.Round((clima.Temp+273)*10) / 10

	return &entity.ClimaOutput{
		City:  loc.Localidade,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}
