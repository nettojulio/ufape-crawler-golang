package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Version string `env:"APP_VERSION" envDefault:"development"`
	Port    int    `env:"APP_PORT" envDefault:"8080"`
	Host    string `env:"APP_HOST" envDefault:"localhost:8080"`
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

	host := "localhost:8080"
	if val, ok := os.LookupEnv("APP_HOST"); ok {
		host = val
	}

	return &Config{
		Version: version,
		Port:    port,
		Host:    host,
	}, nil
}
