package print

import (
	"github.com/mick-roper/rdfox-cli/config"
	"github.com/mick-roper/rdfox-cli/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Cmd() *cobra.Command {
	var cmd cobra.Command
	var path string

	cmd.Use = "print"
	cmd.Short = "prints the current config"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		logger := logging.GetFromContext(cmd.Context())

		logger.Debug("reading file", zap.String("path", path))

		cfg, err := config.File(path)
		if err != nil {
			logger.Error("could not read file", zap.Error(err))
			return err
		}

		logger.Info("got config", zap.Any("config", cfg))

		return nil
	}

	cmd.Flags().StringVar(&path, "path", config.DefaultFilePath, "the path to the config file")

	return &cmd
}
