package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/sudores/invoice-system/pkg/api"
	"github.com/sudores/invoice-system/pkg/api/auth"
	"github.com/sudores/invoice-system/pkg/repo"
)

type Config struct {
	RepoConfig repo.Config
	LogLevel   string `env:"LOG_LEVEL", envDefault:"debug"`
	Jwt        auth.Config
	Api        api.Config
}

// Parse parses the env variables defined in Cnf tags to Cnf struct pointer
func Parse() (*Config, error) {
	cnf := Config{}
	if err := env.Parse(&cnf); err != nil {
		return nil, err
	}
	return &cnf, nil
}
