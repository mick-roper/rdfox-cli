package roles

import (
	"errors"

	v6 "github.com/mick-roper/rdfox-cli/rdfox/v6"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func revokePrivileges() *cobra.Command {
	var cmd cobra.Command

	cmd.Use = "revoke"
	cmd.Short = "revoke privileges from a role"

	var roleToUpdate string
	var datastore string
	var resource string
	var accessTypes string

	cmd.Flags().StringVar(&roleToUpdate, "role-to-update", "", "the name of the role to grant privileges to")
	cmd.Flags().StringVar(&datastore, "datastore", "", "the datastore these privileges apply to")
	cmd.Flags().StringVar(&resource, "resource", "*", "the resource these privileges apply to")
	cmd.Flags().StringVar(&accessTypes, "access-types", "", "the access types this role should have")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := utils.LoggerFromContext(ctx)

		if roleToUpdate == "" {
			logger.Error("arg not set", zap.String("arg", "role-to-update"))
			return errors.New("arg not set")
		}

		if datastore == "" {
			logger.Error("arg not set", zap.String("arg", "datastore"))
			return errors.New("arg not set")
		}

		if accessTypes == "" {
			logger.Error("arg not set", zap.String("arg", "access-types"))
			return errors.New("arg not set")
		}

		logger.Debug("getting root command flags...")

		r := utils.RootCommandFlags(cmd)

		logger.Debug("got root command flags", zap.Any("flags", r))
		logger.Debug("revoking privileges...")

		if err := v6.RevokeDatastorePrivileges(ctx, r.Server, r.Protocol, r.Role, r.Password, roleToUpdate, datastore, resource, accessTypes); err != nil {
			logger.Error("could not update role", zap.Error(err))
			return err
		}

		logger.Debug("privileges have been updated")

		return nil
	}

	return &cmd
}
