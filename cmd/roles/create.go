package roles

import (
	"fmt"

	v6 "github.com/mick-roper/rdfox-cli/rdfox/v6"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func createRole() *cobra.Command {
	var cmd cobra.Command

	cmd.Use = "create"
	cmd.Short = "create a new role"

	var newRoleName string
	var newRolePassword string

	cmd.Flags().StringVar(&newRoleName, "new-role-name", "", "the name of the new role")
	cmd.Flags().StringVar(&newRolePassword, "new-role-password", "", "the password of the new role")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := utils.LoggerFromContext(ctx)

		if newRoleName == "" {
			logger.Error("arg not set", zap.String("arg", "new-role-name"))
			return fmt.Errorf("arg not set")
		}

		if newRolePassword == "" {
			logger.Error("arg not set", zap.String("arg", "new-role-password"))
			return fmt.Errorf("arg not set")
		}

		logger.Debug("getting root command flags...")

		r := utils.RootCommandFlags(cmd)

		logger.Debug("got root command flags", zap.Any("flags", r))
		logger.Debug("creating role...")

		if err := v6.CreateRole(ctx, r.Server, r.Protocol, r.Role, r.Password, newRoleName, newRolePassword); err != nil {
			logger.Error("could not create role", zap.Error(err))
			return err
		}

		logger.Debug("role created")

		fmt.Printf("role created")

		return nil
	}

	return &cmd
}
