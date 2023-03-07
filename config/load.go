package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
)

type Config struct {
	Username string `json:"username,omit-empty"`
	Password string `json:"password,omit-empty"`
}

func Load() (*Config, error) {
	dir, err := homedir.Dir()

	if err != nil {
		return nil, err
	}

	path := fmt.Sprint(dir, "/.rdfox-cli/config.json")

	file, err := os.OpenFile(path, os.O_APPEND, os.ModeAppend)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		if file, err = os.Create(path); err != nil {
			return nil, err
		}
	}

	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
