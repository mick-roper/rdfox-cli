package operation

import (
	importaxioms "github.com/mick-roper/rdfox-cli/cmd/operation/import-axioms"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var cmd cobra.Command

	cmd.Use = "operation"
	cmd.Short = "contains 'operation' subcommands"

	cmd.AddCommand(importaxioms.Cmd())

	return &cmd
}
