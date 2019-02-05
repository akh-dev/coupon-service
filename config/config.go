package config

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

type Config struct {
	DB      DBConf
	Service ServiceConf
}

// DBConf - DB config
type DBConf struct {
	Host string `env:"DB_HOST" envDefault:"localhost"`
	Port string `env:"DB_PORT" envDefault:"27017"`
	Name string `env:"DB_NAME" envDefault:"test"`
}

type ServiceConf struct {
	CtxTimeout int    `env:"CONTEXT_TIMEOUT" envDefault:"10"`
	Port       string `env:"LISTEN_PORT" envDefault:"8080"`
	Debug      bool   `env:"DEBUG" envDefault:"true"`
}

func Get() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(&cfg.DB); err != nil {
		return nil, errors.Wrap(err, "Failed to load DB config")
	}

	if err := env.Parse(&cfg.Service); err != nil {
		return nil, errors.Wrap(err, "Failed to load Service config")
	}

	return cfg, nil
}
