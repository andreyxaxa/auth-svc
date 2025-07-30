package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		HTTP    HTTP
		Log     Log
		PG      PG
		Swagger Swagger
		JWT     JWT
		Webhook Webhook
	}

	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	JWT struct {
		SecretKey string `env:"JWT_SECRET_KEY,required"`
		TTL       int    `env:"JWT_TTL,required"`
	}

	Webhook struct {
		URL string `env:"WEBHOOK_URL,required"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
