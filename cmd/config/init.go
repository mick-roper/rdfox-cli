package config

import (
	"github.com/mick-roper/rdfox-cli/config"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type initCommandConfig struct {
	server        string
	protocol      string
	role          string
	password      string
	logLevel      string
	serverVersion int
}

func (c initCommandConfig) Server() string {
	return c.server
}

func (c initCommandConfig) Protocol() string {
	return c.protocol
}

func (c initCommandConfig) Role() string {
	return c.role
}

func (c initCommandConfig) Password() string {
	return c.password
}

func (c initCommandConfig) LogLevel() string {
	return c.logLevel
}

func (c initCommandConfig) ServerVersion() int {
	return c.serverVersion
}

func initCmd() *cobra.Command {
	var x initCommandConfig
	var path string
	var overwrite bool

	var cmd cobra.Command
	cmd.Use = "init"
	cmd.Short = "initialises the config"

	cmd.Flags().StringVar(&x.server, "server", "", "the name of the server")
	cmd.Flags().StringVar(&x.protocol, "protocol", "https", "the protocol to sue to communicate with the server")
	cmd.Flags().StringVar(&x.role, "role", "", "the role to use to connect to the server")
	cmd.Flags().StringVar(&x.password, "password", "", "the password to use to connect to the server")
	cmd.Flags().StringVar(&x.logLevel, "default-log-level", "info", "the log level to use as default")
	cmd.Flags().StringVar(&path, "path", config.DefaultFilePath(), "the config file path")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "<true> to overwrite an existing config")

	cmd.MarkFlagsRequiredTogether("server", "role", "password")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := utils.LoggerFromContext(ctx)

		logger.Debug("checking required flags...")

		if err := cmd.ValidateRequiredFlags(); err != nil {
			logger.Error("required flags cehck failed", zap.Error(err))
			return err
		}

		logger.Debug("flags are valid - writing the config to file...")

		if err := config.WriteFile(ctx, path, x, overwrite); err != nil {
			logger.Error("could not write config file", zap.Error(err))
			return err
		}

		logger.Debug("config file written!")

		return nil
	}

	return &cmd
}
