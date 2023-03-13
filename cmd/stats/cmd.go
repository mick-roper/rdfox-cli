package stats

import (
	"github.com/mick-roper/rdfox-cli/logging"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var cmd cobra.Command

	cmd.Use = "stats"
	cmd.Short = "get stats for a server or datastore"

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := logging.GetFromContext(ctx)

		logger.Debug("we should write stuff here!")

		return nil
	}

	return &cmd
}
