package operation

import (
	"errors"

	v6 "github.com/mick-roper/rdfox-cli/rdfox/v6"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func importAxiomsCommand() *cobra.Command {
	var cmd cobra.Command
	var datastore string
	var srcGraph string
	var dstGraph string

	cmd.Use = "import-axioms"
	cmd.Short = "imports axioms from a source graph, and imports them into a destination graph"

	cmd.Flags().StringVar(&datastore, "datastore", "", "the datastore")
	cmd.Flags().StringVar(&srcGraph, "src-graph", "", "the source graph")
	cmd.Flags().StringVar(&dstGraph, "dst-graph", "", "the destination graph")

	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		if datastore == "" {
			return errors.New("datastore is unset")
		}

		if srcGraph == "" {
			return errors.New("src-graph is unset")
		}

		if dstGraph == "" {
			return errors.New("dst-graph is unset")
		}

		ctx := cmd.Context()
		logger := utils.LoggerFromContext(ctx)

		logger.Debug("getting flags...")

		server := cmd.Flags().Lookup("server").Value.String()
		protocol := cmd.Flags().Lookup("protocol").Value.String()
		role := cmd.Flags().Lookup("role").Value.String()
		password := cmd.Flags().Lookup("password").Value.String()

		logger.Debug("got flags", zap.String("server", server), zap.String("protocol", protocol), zap.String("role", role), zap.String("password", password))
		logger.Debug("importing axioms...")

		if err := v6.ImportAxioms(ctx, protocol, server, role, password, datastore, srcGraph, dstGraph); err != nil {
			logger.Error("could not import axioms", zap.Error(err))
			return err
		}

		logger.Debug("axioms imported")

		return nil
	}

	return &cmd
}
