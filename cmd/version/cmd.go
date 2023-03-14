package version

import (
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Cmd(currentVersion string) *cobra.Command {
	var cmd cobra.Command

	cmd.Use = "version"

	cmd.Run = func(cmd *cobra.Command, args []string) {
		logger := utils.LoggerFromContext(cmd.Context())

		logger.Info("current version", zap.String("version", currentVersion))
	}

	return &cmd
}
