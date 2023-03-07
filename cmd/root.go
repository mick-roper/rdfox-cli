package cmd

import (
	"context"
	"time"

	"github.com/mick-roper/rdfox-cli/cmd/config"
	"github.com/mick-roper/rdfox-cli/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Execute() int {
	logger := logging.New()
	defer logger.Sync()

	ctx := context.TODO()
	ctx = logging.AddToContext(ctx, logger)
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)

	defer cancel()

	cmd := newRootCommand()
	cmd.AddCommand(config.Cmd(nil))

	if err := cmd.ExecuteContext(ctx); err != nil {
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

	return &cmd
}
