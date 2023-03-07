package config

import (
	"github.com/mick-roper/rdfox-cli/config"
	"github.com/mick-roper/rdfox-cli/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Cmd(cfg *config.Config) *cobra.Command {
	cmd := cobra.Command{
		Use:   "config",
		Short: "shows the current config",
		Long:  "shows the contents of the configuration file",
		Run: func(cmd *cobra.Command, _ []string) {
			logger := logging.GetFromContext(cmd.Context())

			logger.Info("config", zap.Any("contents", cfg))
		},
	}

	return &cmd
}
