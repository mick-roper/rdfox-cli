package config

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var cmd cobra.Command
	cmd.Use = "config"
	cmd.Short = "configures the CLI"

	cmd.AddCommand(printCmd())
	cmd.AddCommand(initCmd())

	return &cmd
}
