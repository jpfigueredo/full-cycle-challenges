package repository_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jpfigueredo/cep-clima-challenge/internal/repository"
)

func TestWeatherAPIRepoSucesso(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			if strings.Contains(req.URL.String(), "api.weatherapi.com") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"current": {"temp_c": 25.5}}`)),
				}
			}
			return &http.Response{StatusCode: http.StatusNotFound, Body: http.NoBody}
		}),
	}

	repo := &repository.WeatherAPIRepo{APIKey: "testkey", Client: client}
	clima, err := repo.GetClimaAtual(context.Background(), "São Paulo")
	if err != nil || clima.Temp != 25.5 {
		t.Errorf("Expected sucesso com temp 25.5, got err=%v, temp=%f", err, clima.Temp)
	}
}

func TestWeatherAPIRepoNotOK(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			if strings.Contains(req.URL.String(), "api.weatherapi.com") {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       http.NoBody,
				}
			}
			return &http.Response{StatusCode: http.StatusNotFound, Body: http.NoBody}
		}),
	}

	repo := &repository.WeatherAPIRepo{APIKey: "testkey", Client: client}
	_, err := repo.GetClimaAtual(context.Background(), "Invalid")
	if err == nil {
		t.Error("Expected erro para status not OK")
	}
}

func TestWeatherAPIRepoInvalidJSON(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			if strings.Contains(req.URL.String(), "api.weatherapi.com") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`invalid`)),
				}
			}
			return &http.Response{StatusCode: http.StatusNotFound, Body: http.NoBody}
		}),
	}

	repo := &repository.WeatherAPIRepo{APIKey: "testkey", Client: client}
	_, err := repo.GetClimaAtual(context.Background(), "São Paulo")
	if err == nil {
		t.Error("Expected erro para JSON inválido")
	}
}
