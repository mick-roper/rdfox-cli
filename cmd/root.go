package cmd

import (
	"context"
	"time"

	"github.com/mick-roper/rdfox-cli/cmd/stats"
	"github.com/mick-roper/rdfox-cli/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Execute() int {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)

	defer cancel()

	cmd := newRootCommand()
	cmd.AddCommand(stats.Cmd())

	preRun := func(cmd *cobra.Command, _ []string) {
		level := cmd.Flags().Lookup("log-level").Value.String()
		logger := logging.New(level)
		ctx = logging.AddToContext(cmd.Context(), logger)
		cmd.SetContext(ctx)
	}

	postRun := func(cmd *cobra.Command, _ []string) {
		logging.GetFromContext(cmd.Context()).Sync()
	}

	cmd.PersistentPreRun = preRun
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		preRun(cmd, args)
		return nil
	}

	cmd.PersistentPostRun = postRun
	cmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		postRun(cmd, args)
		return nil
	}

	if err := cmd.ExecuteContext(ctx); err != nil {
		logger := logging.GetFromContext(cmd.Context())
		logger.Error("execution failed", zap.Error(err))
		return 1
	}

	return 0
}

func newRootCommand() *cobra.Command {
	var cmd cobra.Command

	flags := cmd.PersistentFlags()
	flags.String("log-level", "info", "the log level used by the CLI")
	flags.String("role", "", "the role used to communicate with RDFox")
	flags.String("password", "", "the password used to communicate with RDFox")
	flags.String("server", "", "the name of the RDFox server")
	flags.String("protocol", "https", "the protocol to use to communicate with RDFox")

	return &cmd
}
