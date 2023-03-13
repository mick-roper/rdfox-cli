package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

const DefaultFilePath = "$HOME/.rdfox-cli"

const (
	keyServer   = "server"
	keyProtocol = "protocol"
	keyRole     = "role"
	keyPassword = "password"
	keyLogLevel = "log_level"
)

type fileConfig struct {
	server   string
	protocol string
	role     string
	password string
	logLevel string
}

func (f fileConfig) Server() string {
	return f.server
}

func (f fileConfig) Protocol() string {
	return f.protocol
}

func (f fileConfig) Role() string {
	return f.role
}

func (f fileConfig) Password() string {
	return f.password
}

func (f fileConfig) LogLevel() string {
	return f.logLevel
}

func DefaultFile() (Config, error) {
	return File(DefaultFilePath)
}

func File(path string) (Config, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, errors.New("file does not exist")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var cfg fileConfig

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		s := scanner.Text()
		parts := strings.SplitN(s, "\t", 2)

		if len(parts) != 2 {
			return nil, errors.New("invalid file format")
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case keyServer:
			cfg.server = value
		case keyProtocol:
			cfg.protocol = value
		case keyRole:
			cfg.role = value
		case keyPassword:
			cfg.password = value
		case keyLogLevel:
			cfg.logLevel = value
		}
	}

	return &cfg, nil
}

func WriteFile(cfg Config) error {
	return errors.New("not implemented")
}
