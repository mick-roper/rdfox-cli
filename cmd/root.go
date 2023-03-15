package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mick-roper/rdfox-cli/cmd/config"
	exportdata "github.com/mick-roper/rdfox-cli/cmd/export-data"
	"github.com/mick-roper/rdfox-cli/cmd/operation"
	"github.com/mick-roper/rdfox-cli/cmd/stats"
	"github.com/mick-roper/rdfox-cli/cmd/version"
	configuration "github.com/mick-roper/rdfox-cli/config"
	"github.com/mick-roper/rdfox-cli/logging"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Execute(currentVersion string) int {
	ctx, cancel := context.WithCancel(context.TODO())

	defer cancel()

	ctx = utils.AddHttpClientToContext(ctx, http.DefaultClient)

	cmd := newRootCommand(ctx)
	cmd.AddCommand(version.Cmd(currentVersion))
	cmd.AddCommand(stats.Cmd())
	cmd.AddCommand(config.Cmd())
	cmd.AddCommand(operation.Cmd())
	cmd.AddCommand(exportdata.Cmd())

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

	okChan := make(chan struct{})
	errChan := make(chan error)
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	defer close(okChan)
	defer close(errChan)
	defer close(sigChan)

	go func() {
		if err := cmd.ExecuteContext(ctx); err != nil {
			errChan <- err
			return
		}

		okChan <- struct{}{}
	}()

	var exitCode int

	select {
	case <-okChan:
		exitCode = 0
	case <-sigChan:
		exitCode = 0
	case err := <-errChan:
		utils.LoggerFromContext(ctx).Error("execution failed", zap.Error(err))
		exitCode = 1
	}

	return exitCode
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
