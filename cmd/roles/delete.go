package roles

import (
	"errors"
	"fmt"

	"github.com/mick-roper/rdfox-cli/console"
	"github.com/mick-roper/rdfox-cli/rdfox"
	v6 "github.com/mick-roper/rdfox-cli/rdfox/v6"
	v7 "github.com/mick-roper/rdfox-cli/rdfox/v7"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func deleteRole() *cobra.Command {
	var cmd cobra.Command

	cmd.Use = "delete"
	cmd.Short = "deletes a role"

	var roleToDelete string

	cmd.Flags().StringVar(&roleToDelete, "role-to-delete", "", "the name of the role that should be deleted")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := utils.LoggerFromContext(ctx)

		if roleToDelete == "" {
			logger.Error("arg not set", zap.String("arg", "role-to-delete"))
			return errors.New("arg not set")
		}

		logger.Debug("asking for confirmation...")

		if ok := console.BoolPrompt("are you sure you want to delete the role?"); !ok {
			logger.Info("you must provide confirmation that you want to delete the role")
			return nil
		}

		logger.Debug("got confirmation")
		logger.Debug("getting root command flags...")

		r := utils.RootCommandFlags(cmd)
		var fn rdfox.DeleteRole
		switch r.Version {
		case 6:
			fn = v6.DeleteRole
		case 7:
			fn = v7.DeleteRole
		default:
			return fmt.Errorf("RDFox version %d is unsupported", r.Version)
		}

		logger.Debug("got root command flags", zap.Any("flags", r))
		logger.Debug("deleting role...")

		if err := fn(ctx, r.Server, r.Protocol, r.Role, r.Password, roleToDelete); err != nil {
			logger.Error("could not delete role", zap.Error(err))
			return err
		}

		logger.Info("role deleted")

		return nil
	}

	return &cmd
}
