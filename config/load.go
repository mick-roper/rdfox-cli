package config

import "errors"

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func Load() (*Config, error) {
	return nil, errors.New("not implemented")
}
