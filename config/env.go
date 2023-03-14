package config

import "os"

type envConfig struct{}

func (envConfig) Server() string {
	return os.Getenv("RDFOX_CLI_SERVER")
}

func (envConfig) Protocol() string {
	if s := os.Getenv("RDFOX_CLI_PROTOCOL"); s != "" {
		return s
	}

	return "https"
}

func (envConfig) Role() string {
	return os.Getenv("RDFOX_CLI_ROLE")
}

func (envConfig) Password() string {
	return os.Getenv("RDFOX_CLI_PASSWORD")
}

func (envConfig) LogLevel() string {
	if s := os.Getenv("RDFOX_CLI_LOG_LEVEL"); s != "" {
		return s
	}

	return "info"
}

func FromEnv() Config {
	return envConfig{}
}
