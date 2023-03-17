package roles

import (
	"errors"

	v6 "github.com/mick-roper/rdfox-cli/rdfox/v6"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func getInfo() *cobra.Command {
	var cmd cobra.Command

	cmd.Use = "info"
	cmd.Short = "get info about a role"

	var roleToInspect string

	cmd.Flags().StringVar(&roleToInspect, "role-to-inspect", "", "the name of the role to inspect")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := utils.LoggerFromContext(ctx)

		if roleToInspect == "" {
			logger.Error("arg not set", zap.String("arg", "role-to-inspect"))
			return errors.New("arg not set")
		}

		logger.Debug("getting root command flags...")

		r := utils.RootCommandFlags(cmd)

		logger.Debug("got root command flags", zap.Any("flags", r))
		logger.Debug("getting privileges...")

		privileges, err := v6.ListPrivileges(ctx, r.Server, r.Protocol, r.Role, r.Password, roleToInspect)
		if err != nil {
			logger.Error("could not list privileges", zap.Error(err))
			return err
		}

		logger.Debug("got privileges")

		for resource, accessTypes := range privileges {
			logger.Info("got data", zap.String("resource", resource), zap.Any("access-types", accessTypes))
		}

		return nil
	}

	return &cmd
}
