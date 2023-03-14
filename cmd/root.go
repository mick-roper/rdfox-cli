package cmd

import (
	"context"
	"net/http"
	"time"

	"github.com/mick-roper/rdfox-cli/cmd/config"
	"github.com/mick-roper/rdfox-cli/cmd/stats"
	configuration "github.com/mick-roper/rdfox-cli/config"
	"github.com/mick-roper/rdfox-cli/logging"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Execute() int {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)

	defer cancel()

	ctx = utils.AddHttpClientToContext(ctx, http.DefaultClient)

	cmd := newRootCommand(ctx)
	cmd.AddCommand(stats.Cmd())
	cmd.AddCommand(config.Cmd())

	preRun := func(cmd *cobra.Command, _ []string) {
		level := cmd.Flags().Lookup("log-level").Value.String()
		logger := logging.New(level)
		ctx = utils.AddLoggerToContext(cmd.Context(), logger)
		cmd.SetContext(ctx)
	}

	postRun := func(cmd *cobra.Command, _ []string) {
		utils.LoggerFromContext(cmd.Context()).Sync()
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
		utils.LoggerFromContext(cmd.Context()).Error("execution failed", zap.Error(err))
		return 1
	}

	return 0
}

func newRootCommand(ctx context.Context) *cobra.Command {
	var (
		defaultLogLevel = "info"
		defaultProtocol = "https"
	)

	cfg := configuration.Load(configuration.FileLoader(ctx, configuration.DefaultFilePath()), configuration.FromEnv)

	if s := cfg.LogLevel(); s != "" {
		defaultLogLevel = s
	}

	if s := cfg.Protocol(); s != "" {
		defaultProtocol = s
	}

	defaultServer := cfg.Server()
	defaultRole := cfg.Role()
	defaultPassword := cfg.Password()

	var cmd cobra.Command
	cmd.SetContext(ctx)
	flags := cmd.PersistentFlags()
	flags.String("log-level", defaultLogLevel, "the log level used by the CLI")
	flags.String("role", defaultRole, "the role used to communicate with RDFox")
	flags.String("password", defaultPassword, "the password used to communicate with RDFox")
	flags.String("server", defaultServer, "the name of the RDFox server")
	flags.String("protocol", defaultProtocol, "the protocol to use to communicate with RDFox")

	return &cmd
}
