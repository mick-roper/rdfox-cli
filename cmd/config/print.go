package config

import (
	"github.com/mick-roper/rdfox-cli/config"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func printCmd() *cobra.Command {
	var cmd cobra.Command
	var path string

	cmd.Use = "print"
	cmd.Short = "prints the current config"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := utils.LoggerFromContext(ctx)

		logger.Debug("getting flags...")

		server := cmd.Flags().Lookup("server").Value.String()
		protocol := cmd.Flags().Lookup("protocol").Value.String()
		role := cmd.Flags().Lookup("role").Value.String()
		password := cmd.Flags().Lookup("password").Value.String()
		logLevel := cmd.Flags().Lookup("log-level").Value.String()

		logger.Info("got config",
			zap.String("server", server),
			zap.String("role", role),
			zap.String("password", password),
			zap.String("protocol", protocol),
			zap.String("log-level", logLevel),
		)

		return nil
	}

	cmd.Flags().StringVar(&path, "path", config.DefaultFilePath(), "the path to the config file")

	return &cmd
}
