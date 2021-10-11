package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BackendTable string `envconfig:"BACKEND_TABLE"`
}

func (c *Config) validate() error {
	return nil
}

func NewConfig() *Config {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}
	err = cfg.validate()
	if err != nil {
		panic(err)
	}
	return &cfg
}
