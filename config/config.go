package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Env string

const (
	Env_Test Env = "test"
	Env_Dev  Env = "dev"
)


//struct to the environment varibles.
type Config struct {
	DatabaseName     string  `env:"DATABASE_NAME"`
	DatabaseHost     string  `env:"DATABASE_HOST"`
	DatabaseUser     string  `env:"DATABASE_USER"`
	DatabasePassword string  `env:"DATABASE_PASSWORD"`
	DatabasePort	 string  `env:"DATABASE_PORT"`
	DatabasePortTest string  `env:"DATABASE_TEST_PORT"`
	Env 			 Env     `env:"ENV" envDefault:"dev"`

}

func (c *Config) DatabaseUrl() string {
	port := c.DatabasePort
	if c.Env == Env_Test {
		port = c.DatabasePortTest
	}
	
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		 c.DatabaseUser,
		 c.DatabasePassword,
		 c.DatabaseHost,
		 port,
		 c.DatabaseName,
	)
}

func New() (*Config, error){
	cfg, err := env.ParseAs[Config]();
	if err != nil {
		return nil, fmt.Errorf("Faied to load config: %w", err)
	}
	return &cfg, nil
}