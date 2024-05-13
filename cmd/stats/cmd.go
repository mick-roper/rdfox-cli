package stats

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mick-roper/rdfox-cli/rdfox"
	v6 "github.com/mick-roper/rdfox-cli/rdfox/v6"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type formatter func(rdfox.Statistics) error

func Cmd() *cobra.Command {
	var cmd cobra.Command
	var datastore string
	var format string

	cmd.Use = "stats"
	cmd.Short = "get stats for a server or datastore"

	cmd.Flags().StringVar(&datastore, "datastore", "", "The datastore that you want stats for. Leave blank to get server stats.")
	cmd.Flags().StringVar(&format, "format", "console", "The format of the results (console, json).")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := utils.LoggerFromContext(ctx)

		var f formatter = consoleFormatter

		if format == "json" {
			f = jsonFormatter
		}

		logger.Debug("getting flags...")

		server := cmd.Flags().Lookup("server").Value.String()
		protocol := cmd.Flags().Lookup("protocol").Value.String()
		role := cmd.Flags().Lookup("role").Value.String()
		password := cmd.Flags().Lookup("password").Value.String()

		logger.Debug("got flags", zap.String("server", server), zap.String("protocol", protocol), zap.String("role", role), zap.String("password", password))
		logger.Debug("getting stats...")

		stats, err := v6.GetStats(ctx, server, protocol, role, password, datastore)
		if err != nil {
			logger.Error("could not get stats", zap.Error(err))
			return err
		}

		logger.Debug("got stats", zap.Any("stats", stats))

		if err := f(stats); err != nil {
			logger.Error("could not print stats", zap.Error(err))
			return err
		}

		return nil
	}

	return &cmd
}

func consoleFormatter(s rdfox.Statistics) error {
	for subject, duples := range s {
		fmt.Print(subject)

		for predicate, object := range duples {
			fmt.Print("\n\t", predicate, ":\t", object)
		}

		fmt.Print("\n")
	}

	return nil
}

func jsonFormatter(s rdfox.Statistics) error {
	return json.NewEncoder(os.Stdout).Encode(s)
}
