package roles

import (
	v6 "github.com/mick-roper/rdfox-cli/rdfox/v6"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func listRoles() *cobra.Command {
	var cmd cobra.Command

	cmd.Use = "list"
	cmd.Short = "lists all roles"

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := utils.LoggerFromContext(ctx)

		logger.Debug("getting root flags...")

		r := utils.RootCommandFlags(cmd)

		logger.Debug("got root command flags", zap.Any("flags", r))
		logger.Debug("getting roles...")

		roles, err := v6.GetRoles(ctx, r.Server, r.Protocol, r.Role, r.Password)
		if err != nil {
			logger.Error("could not get roles", zap.Error(err))
			return err
		}

		logger.Debug("got roles", zap.Int("count", len(roles)))

		logger.Info("got roles", zap.Any("roles", roles))

		return nil
	}

	return &cmd
}
