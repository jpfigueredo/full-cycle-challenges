package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BidResponse struct {
	Bid string `json:"bid"`
}

func main() {
	resp, err := http.Get("http://localhost:8080/cotacao")
	if err != nil {
		fmt.Println("Erro ao chamar servidor:", err)
		return
	}
	defer resp.Body.Close()

	var r BidResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		fmt.Println("Erro ao ler resposta:", err)
		return
	}

	fmt.Println("Dólar:", r.Bid)
}
