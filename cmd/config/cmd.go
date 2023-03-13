package config

import (
	"github.com/mick-roper/rdfox-cli/cmd/config/print"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var cmd cobra.Command
	cmd.Use = "config"
	cmd.Short = "configures the CLI"

	cmd.AddCommand(print.Cmd())

	return &cmd
}
