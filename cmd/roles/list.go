package roles

import (
	"fmt"

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

		logger.Debug("getting roles...")

		r := utils.RootCommandFlags(cmd)
		roles, err := v6.GetRoles(ctx, r.Server, r.Protocol, r.Role, r.Password)
		if err != nil {
			logger.Error("coudl not get roles", zap.Error(err))
			return err
		}

		logger.Debug("got roles", zap.Int("count", len(roles)))

		for _, role := range roles {
			fmt.Println(role)
		}

		return nil
	}

	return &cmd
}
