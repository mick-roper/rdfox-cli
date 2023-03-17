package roles

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	var cmd cobra.Command

	cmd.Use = "roles"
	cmd.Short = "manage roles"
	cmd.Long = "provides role management functionality"

	cmd.AddCommand(listRoles())

	return &cmd
}
