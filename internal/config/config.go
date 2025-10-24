package config

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Version string `env:"APP_VERSION"`
	Port    int    `env:"APP_PORT" envDefault:"8080"`
	Host    string `env:"APP_HOST" envDefault:"localhost:8080"`
}

// Load carrega as configurações da aplicação
func Load(version string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using environment variables")
	}

	cfg := Config{
		Version: version,
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from environment: %w", err)
	}

	return &cfg, nil
}
