package operation

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var cmd cobra.Command

	cmd.Use = "operation"
	cmd.Short = "contains 'operation' subcommands"

	cmd.AddCommand(importAxiomsCommand())

	return &cmd
}
