package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

//struct to the environment varibles.
type Config struct {
	DatabaseName     string  `env:"DATABASE_NAME"`
	DatabaseHost     string  `env:"DATABASE_HOST"`
	DatabaseUser     string  `env:"DATABASE_USER"`
	DatabasePassword string  `env:"DATABASE_PASSWORD"`
	DatabasePort	 string  `env:"DATABASE_PORT"`

}

func New() (*Config, error){
	cfg, err := env.ParseAs[Config]();
	if err != nil {
		return nil, fmt.Errorf("Faied to load config: %w", err)
	}
	return &cfg, nil
}