package entity

type Localizacao struct {
	Localidade string `json:"localidade"`
}

type ClimaInput struct {
	Temp float64 `json:"temp"`
}

type ClimaOutput struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}
