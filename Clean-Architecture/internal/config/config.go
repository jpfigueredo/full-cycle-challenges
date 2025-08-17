package config

import (
	"os"
)

type Config struct {
	ServerPort string
}

func Load() *Config {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		ServerPort: port,
	}
}
