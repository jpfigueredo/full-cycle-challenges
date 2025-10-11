package repository_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jpfigueredo/cep-clima-distributed/internal/shared/repository"
)

type roundTripFunc func(*http.Request) *http.Response

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r), nil }

func TestViaCEPRepoSucesso(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			if strings.Contains(req.URL.String(), "viacep.com.br") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"localidade": "São Paulo"}`)),
				}
			}
			return &http.Response{StatusCode: http.StatusNotFound, Body: http.NoBody}
		}),
	}

	repo := &repository.ViaCEPRepo{Client: client}
	loc, err := repo.GetLocalizacao(context.Background(), "98765432")
	if err != nil || loc.Localidade != "São Paulo" {
		t.Errorf("Expected sucesso com localidade São Paulo, got err=%v, localidade=%s", err, loc.Localidade)
	}
}

func TestViaCEPRepoNotFound(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			if strings.Contains(req.URL.String(), "viacep.com.br") {
				return &http.Response{
					StatusCode: http.StatusOK, // ViaCEP retorna 200 mesmo para not found, com {}
					Body:       io.NopCloser(strings.NewReader(`{}`)),
				}
			}
			return &http.Response{StatusCode: http.StatusNotFound, Body: http.NoBody}
		}),
	}

	repo := &repository.ViaCEPRepo{Client: client}
	loc, err := repo.GetLocalizacao(context.Background(), "00000000")
	if err != nil || loc.Localidade != "" {
		t.Errorf("Expected vazio para CEP not found, got err=%v, localidade=%s", err, loc.Localidade)
	}
}

func TestViaCEPRepoInvalidJSON(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			if strings.Contains(req.URL.String(), "viacep.com.br") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`invalid json`)),
				}
			}
			return &http.Response{StatusCode: http.StatusNotFound, Body: http.NoBody}
		}),
	}

	repo := &repository.ViaCEPRepo{Client: client}
	_, err := repo.GetLocalizacao(context.Background(), "12345678")
	if err == nil {
		t.Error("Expected erro para JSON inválido, got nil")
	}
}
