package config

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mick-roper/rdfox-cli/logging"
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

	var cfg fileConfig

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

func WriteFile(ctx context.Context, path string, cfg Config, overwrite bool) error {
	logger := logging.GetFromContext(ctx).With(zap.String("path", path))

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

	writeFn := func(pairs map[string]string) error {
		for k, v := range pairs {
			s := fmt.Sprintf("%s%s%s\n", k, separator, v)
			if _, err := file.WriteString(s); err != nil {
				return err
			}
		}

		return nil
	}

	logger.Debug("writing file contents")

	data := map[string]string{
		keyServer:   cfg.Server(),
		keyProtocol: cfg.Protocol(),
		keyRole:     cfg.Role(),
		keyPassword: cfg.Password(),
		keyLogLevel: cfg.LogLevel(),
	}

	if err := writeFn(data); err != nil {
		logger.Error("coudl not write file data", zap.Error(err))
		return err
	}

	logger.Debug("file written")

	return nil
}
