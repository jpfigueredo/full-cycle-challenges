package repository

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/jpfigueredo/cep-clima-distributed/service-a/internal/entity"
)

type WeatherAPIRepo struct {
	APIKey string
	Client *http.Client
}

func NewWeatherAPIRepo() *WeatherAPIRepo {
	key := os.Getenv("WEATHER_API_KEY")
	return &WeatherAPIRepo{APIKey: key, Client: http.DefaultClient}
}

func (r *WeatherAPIRepo) GetClimaAtual(ctx context.Context, localidade string) (*entity.ClimaInput, error) {
	url := "http://api.weatherapi.com/v1/current.json?key=" + r.APIKey + "&q=" + url.QueryEscape(localidade)
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
		return nil, http.ErrNoLocation
	}

	var data struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &entity.ClimaInput{Temp: data.Current.TempC}, nil
}
