package stats

import (
	v6 "github.com/mick-roper/rdfox-cli/rdfox/v6"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Cmd() *cobra.Command {
	var cmd cobra.Command
	var datastore string

	cmd.Use = "stats"
	cmd.Short = "get stats for a server or datastore"

	cmd.Flags().StringVar(&datastore, "datastore", "", "The datastore that you want stats for. Leave blank to get server stats.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := utils.LoggerFromContext(ctx)

		logger.Debug("getting flags...")

		server := cmd.Flags().Lookup("server").Value.String()
		protocol := cmd.Flags().Lookup("protocol").Value.String()
		role := cmd.Flags().Lookup("role").Value.String()
		password := cmd.Flags().Lookup("password").Value.String()

		logger.Debug("got flags", zap.String("server", server), zap.String("protocol", protocol), zap.String("role", role), zap.String("password", password))
		logger.Debug("getting stats...")

		stats, err := v6.GetStats(ctx, server, protocol, role, password, datastore)
		if err != nil {
			logger.Error("could not get stats", zap.Error(err))
			return err
		}

		logger.Info("got stats", zap.Any("stats", stats))

		return nil
	}

	return &cmd
}
