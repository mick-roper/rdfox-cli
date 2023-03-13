package cmd

import (
	"context"
	"time"

	"github.com/mick-roper/rdfox-cli/cmd/config"
	"github.com/mick-roper/rdfox-cli/cmd/stats"
	configuration "github.com/mick-roper/rdfox-cli/config"
	"github.com/mick-roper/rdfox-cli/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Execute() int {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)

	defer cancel()

	cmd := newRootCommand(ctx)
	cmd.AddCommand(stats.Cmd())
	cmd.AddCommand(config.Cmd())

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

func newRootCommand(ctx context.Context) *cobra.Command {
	var (
		defaultLogLevel = "info"
		defaultProtocol = "https"
		defaultRole     = ""
		defaultPassword = ""
		defaultServer   = ""
	)
	var cmd cobra.Command
	cmd.SetContext(ctx)
	cfg, err := configuration.DefaultFile(ctx)
	if err == nil {
		if s := cfg.LogLevel(); s != "" {
			defaultLogLevel = s
		}

		if s := cfg.Protocol(); s != "" {
			defaultProtocol = s
		}

		defaultRole = cfg.Role()
		defaultPassword = cfg.Password()
		defaultServer = cfg.Server()
	}

	flags := cmd.PersistentFlags()
	flags.String("log-level", defaultLogLevel, "the log level used by the CLI")
	flags.String("role", defaultRole, "the role used to communicate with RDFox")
	flags.String("password", defaultPassword, "the password used to communicate with RDFox")
	flags.String("server", defaultServer, "the name of the RDFox server")
	flags.String("protocol", defaultProtocol, "the protocol to use to communicate with RDFox")

	return &cmd
}
