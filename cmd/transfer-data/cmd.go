package exportdata

import (
	"errors"
	"fmt"
	"strings"

	v7 "github.com/mick-roper/rdfox-cli/rdfox/v7"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Cmd() *cobra.Command {
	var (
		cmd         cobra.Command
		datastore   string
		limit       int
		graph       string
		export      string
		destination string
	)

	cmd.Use = "transfer-data"
	cmd.Short = "transfer data from one RDFox instance to another RDFox instance"
	cmd.Long = "TODO: write sometihn inspiring here!"

	cmd.Flags().StringVar(&datastore, "datastore", "", "the datastore that contains the data you want to export")
	cmd.Flags().IntVar(&limit, "limit", 5000, "the maximum number of triples to return in a single cursor request")
	cmd.Flags().StringVar(&graph, "graph", "", "the graph that contains the data you want to export")
	cmd.Flags().StringVar(&export, "export", "all", "the types of facts to export: options are 'all', 'explicit' or 'implicit'")
	cmd.Flags().StringVar(&destination, "dst", "", "the destination server")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if datastore == "" {
			return errors.New("datastore is unset")
		}

		if graph == "" {
			return errors.New("graph is unset")
		}

		if destination == "" {
			return errors.New("destination is unset")
		}

		graph = strings.TrimPrefix(graph, "<")
		graph = strings.TrimSuffix(graph, ">")

		ctx := cmd.Context()
		logger := utils.LoggerFromContext(ctx)

		logger.Debug("getting flags...")

		server := cmd.Flags().Lookup("server").Value.String()
		protocol := cmd.Flags().Lookup("protocol").Value.String()
		role := cmd.Flags().Lookup("role").Value.String()
		password := cmd.Flags().Lookup("password").Value.String()

		logger.Debug("got flags", zap.String("server", server), zap.String("protocol", protocol), zap.String("role", role), zap.String("password", password))
		logger.Debug("creating a connection...")

		connectionID, err := v7.CreateConnection(ctx, server, protocol, role, password, datastore)
		if err != nil {
			logger.Error("could not create a connection", zap.Error(err))
			return err
		}

		defer func() {
			logger.Debug("deleting the connection...")

			if err := v7.DeleteConnection(ctx, server, protocol, role, password, datastore, connectionID); err != nil {
				logger.Error("could not delete connection", zap.Error(err))
			}

			logger.Debug("connection deleted!")
		}()

		logger.Debug("connection created", zap.String("connection-id", connectionID))

		if err := v7.CreateReadOnlyTransaction(ctx, server, protocol, role, password, datastore, connectionID); err != nil {
			logger.Error("coudl not create a transaction", zap.Error(err))
			return err
		}

		defer func() {
			if err := v7.RollbackTransaction(ctx, server, protocol, role, password, datastore, connectionID); err != nil {
				logger.Error("coudl not roll back the transaction", zap.Error(err))
			}
		}()

		logger.Debug("building query...")
		var query string

		switch export {
		case "all":
			query = fmt.Sprintf("SELECT ?s ?p ?o FROM <%s> WHERE { ?s ?p ?o }", graph)
		case "explicit":
			query = fmt.Sprintf("SELECT ?s ?p ?o FROM <%s> WHERE { ?s ?p ?o EXPLICIT TRUE }", graph)
		case "implicit":
			query = fmt.Sprintf("SELECT ?s ?p ?o FROM <%s> WHERE { ?s ?p ?o EXPLICIT FALSE }", graph)
		default:
			return errors.New("export must be one of 'all', 'explicit' or 'implicit'")
		}

		logger.Debug("query built", zap.String("query", query))

		logger.Debug("creating a cursor...")

		cursor, err := v7.CreateTripleCursor(ctx, server, protocol, role, password, datastore, connectionID, query, limit)
		if err != nil {
			logger.Error("could not create a cursor", zap.Error(err))
			return err
		}

		defer func() {
			logger.Debug("closing cursor...")

			if err := cursor.Close(ctx); err != nil {
				logger.Error("could not close the cursor", zap.Error(err))
			}

			logger.Debug("cursor closed!")
		}()

		logger.Debug("cursor created")

		logger.Info("getting data...")

		dataChan := make(chan map[string]map[string][]string)
		readDoneChan := make(chan struct{})
		writeDoneChan := make(chan struct{})

		readData := func() error {
			defer close(readDoneChan)
			defer close(dataChan)

			cursor := cursor.(*v7.TripleCursor)

			for cursor.Read(ctx) {
				if err := cursor.Err; err != nil {
					logger.Error("could not read data", zap.Error(err))
					return err
				}

				dataChan <- cursor.Data
			}

			return nil
		}

		if err := utils.DoWithTicker(readData, func() {
			logger.Info("still getting data...")
		}); err != nil {
			return err
		}
		<-writeDoneChan
		return nil
	}

	return &cmd
}
