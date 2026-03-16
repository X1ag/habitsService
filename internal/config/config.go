package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	DB PostgresConfig	
}

type PostgresConfig	struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"postgres"`
	DSN   string `env:"DB_DSN" envDefault:"postgresql://postgres:postgres@db:5432/postgres?sslmode=disable"`
}

func Load() (*Config, error) {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    return cfg, nil
}

