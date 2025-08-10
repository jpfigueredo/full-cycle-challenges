package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// AddressFetcher é a interface para qualquer fonte de CEP.
type AddressFetcher interface {
	Fetch(ctx context.Context, cep string) (Address, error)
}

// brasilAPI implementa AddressFetcher.
type brasilAPI struct{}

func (b brasilAPI) Fetch(ctx context.Context, cep string) (Address, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Address{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Address{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Address{}, fmt.Errorf("brasilapi: status %d", resp.StatusCode)
	}

	var bResp struct {
		Cep          string `json:"cep"`
		State        string `json:"state"`
		City         string `json:"city"`
		Neighborhood string `json:"neighborhood"`
		Street       string `json:"street"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&bResp); err != nil {
		return Address{}, err
	}

	return Address{
		Cep:          bResp.Cep,
		Street:       bResp.Street,
		Neighborhood: bResp.Neighborhood,
		City:         bResp.City,
		State:        bResp.State,
		Source:       "brasilapi",
	}, nil
}

// viaCEP implementa AddressFetcher.
type viaCEP struct{}

func (v viaCEP) Fetch(ctx context.Context, cep string) (Address, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Address{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Address{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Address{}, fmt.Errorf("viacep: status %d", resp.StatusCode)
	}

	var vResp struct {
		Cep        string `json:"cep"`
		Logradouro string `json:"logradouro"`
		Bairro     string `json:"bairro"`
		Localidade string `json:"localidade"`
		Uf         string `json:"uf"`
		Erro       bool   `json:"erro,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&vResp); err != nil {
		return Address{}, err
	}

	if vResp.Erro {
		return Address{}, errors.New("viacep: cep not found")
	}

	return Address{
		Cep:          vResp.Cep,
		Street:       vResp.Logradouro,
		Neighborhood: vResp.Bairro,
		City:         vResp.Localidade,
		State:        vResp.Uf,
		Source:       "viacep",
	}, nil
}
