package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Version string
	Port    int
}

func Load(version string) (*Config, error) {
	_ = godotenv.Load()

	port := 8080
	if val, ok := os.LookupEnv("APP_PORT"); ok {
		if p, err := strconv.Atoi(val); err == nil {
			port = p
		} else {
			log.Printf("Invalid port in .env: %s, using default %d", val, port)
		}
	}

	return &Config{
		Version: version,
		Port:    port,
	}, nil
}
