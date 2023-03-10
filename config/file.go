package config

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mick-roper/rdfox-cli/utils"
	"go.uber.org/zap"
)

const (
	keyServer   = "server"
	keyProtocol = "protocol"
	keyRole     = "role"
	keyPassword = "password"
	keyLogLevel = "log_level"
)

const separator = "\t"

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

func DefaultFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return fmt.Sprint(home, "/.rdfox-cli")
}

func DefaultFile(ctx context.Context) (Config, error) {
	return File(ctx, DefaultFilePath())
}

func FileLoader(ctx context.Context, path string) loader {
	return func() Config {
		cfg, err := File(ctx, path)
		if err != nil {
			panic(err)
		}

		return cfg
	}
}

func File(ctx context.Context, path string) (Config, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, errors.New("file does not exist")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := read(file)
	if err != nil {
		return nil, err
	}

	var cfg fileConfig
	for k, v := range data {
		switch k {
		case keyLogLevel:
			cfg.logLevel = v
		case keyPassword:
			cfg.password = v
		case keyProtocol:
			cfg.protocol = v
		case keyRole:
			cfg.role = v
		case keyServer:
			cfg.server = v
		}
	}

	return &cfg, nil
}

func WriteFile(ctx context.Context, path string, cfg Config, overwrite bool) error {
	logger := utils.LoggerFromContext(ctx).With(zap.String("path", path))

	var file *os.File

	logger.Debug("getting file stats...")

	_, err := os.Stat(path)

	if err == nil {
		if !overwrite {
			return errors.New("file exists but 'overwrite' is <false> - delete the file manually or set 'overwrite' to <true>")
		}

		logger.Debug("file exists - opening the file...")

		file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			logger.Error("could not open file", zap.Error(err))
			return err
		}

		logger.Debug("file opened - truncating the file...")

		if err := file.Truncate(0); err != nil {
			logger.Error("could not truncate the file", zap.Error(err))
			return err
		}

		logger.Debug("file truncated")
	} else {
		if os.IsNotExist(err) {
			logger.Debug("file does not exist - creating a new file")

			file, err = os.Create(path)
			if err != nil {
				logger.Error("could not create the file", zap.Error(err))
				return err
			}

			logger.Debug("file created")
		} else {
			logger.Error("could not get file stats", zap.Error(err))
			return err
		}

	}

	defer file.Close()

	logger.Debug("writing file contents")

	data := map[string]string{
		keyServer:   cfg.Server(),
		keyProtocol: cfg.Protocol(),
		keyRole:     cfg.Role(),
		keyPassword: cfg.Password(),
		keyLogLevel: cfg.LogLevel(),
	}

	if err := write(file, data); err != nil {
		logger.Error("coudl not write file data", zap.Error(err))
		return err
	}

	logger.Debug("file written")

	return nil
}

func read(file *os.File) (map[string]string, error) {
	m := map[string]string{}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		s := scanner.Text()
		parts := strings.SplitN(s, separator, 2)

		if len(parts) != 2 {
			return nil, errors.New("invalid file format")
		}

		key := parts[0]
		value := parts[1]
		m[key] = value
	}

	return m, nil
}

func write(file *os.File, pairs map[string]string) error {
	for k, v := range pairs {
		s := fmt.Sprintf("%s%s%s\n", k, separator, v)
		if _, err := file.WriteString(s); err != nil {
			return err
		}
	}

	return nil
}
