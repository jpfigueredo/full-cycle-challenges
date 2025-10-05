package repository

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jpfigueredo/cep-clima-distributed/service-a/internal/entity"
)

type ViaCEPRepo struct {
	Client *http.Client
}

func NewViaCEPRepo() *ViaCEPRepo {
	return &ViaCEPRepo{Client: http.DefaultClient}
}

func (r *ViaCEPRepo) GetLocalizacao(ctx context.Context, cep string) (*entity.Localizacao, error) {
	url := "https://viacep.com.br/ws/" + cep + "/json/"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &entity.Localizacao{}, nil
	}

	var loc entity.Localizacao
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		return nil, err
	}

	return &loc, nil
}
