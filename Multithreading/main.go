package main

import (
	"context"
	"fmt"
	"os"
	"strings"
)

func printAddress(addr Address) {
	fmt.Printf("Source: %s\n", addr.Source)
	fmt.Printf("CEP: %s\n", addr.Cep)
	fmt.Printf("Street: %s\n", addr.Street)
	fmt.Printf("Neighborhood: %s\n", addr.Neighborhood)
	fmt.Printf("City: %s\n", addr.City)
	fmt.Printf("State: %s\n", addr.State)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run main.go <cep>")
		os.Exit(1)
	}
	cep := strings.TrimSpace(os.Args[1])
	if cep == "" {
		fmt.Println("empty cep")
		os.Exit(1)
	}

	fetchers := []AddressFetcher{
		brasilAPI{},
		viaCEP{},
	}
	raceService := NewRaceService(fetchers, 1_000_000_000)

	addr, err := raceService.Run(context.Background(), cep)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	printAddress(addr)
}
