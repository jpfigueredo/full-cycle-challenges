package main

import (
	"log"

	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/api"

	"github.com/jpfigueredo/full-cycle-challenges/Clean-Architecture/internal/config"
)

func main() {
	cfg := config.Load()

	r := api.SetupRouter()

	log.Printf("🚀 Server running on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("❌ failed to start server: %v", err)
	}
}
