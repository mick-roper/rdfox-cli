package exportdata

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	v6 "github.com/mick-roper/rdfox-cli/rdfox/v6"
	"github.com/mick-roper/rdfox-cli/ttl"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Cmd() *cobra.Command {
	var cmd cobra.Command

	var datastore string
	var filePath string
	var limit int
	var graph string

	cmd.Use = "export-data"
	cmd.Short = "export data from the database"
	cmd.Long = "TODO: write sometihn inspiring here!"

	cmd.Flags().StringVar(&datastore, "datastore", "", "the datastore that contains the data you want to export")
	cmd.Flags().StringVar(&filePath, "file", "export.ttl", "the file that the exported data will be written to")
	cmd.Flags().IntVar(&limit, "limit", 5000, "the maximum number of triples to return in a single cursor request")
	cmd.Flags().StringVar(&graph, "graph", "", "the graph that contains the data you want to export")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if datastore == "" {
			return errors.New("datastore is unset")
		}

		if graph == "" {
			return errors.New("graph is unset")
		}

		ctx := cmd.Context()
		logger := utils.LoggerFromContext(ctx)

		logger.Debug("getting flags...")

		server := cmd.Flags().Lookup("server").Value.String()
		protocol := cmd.Flags().Lookup("protocol").Value.String()
		role := cmd.Flags().Lookup("role").Value.String()
		password := cmd.Flags().Lookup("password").Value.String()

		logger.Debug("got flags", zap.String("server", server), zap.String("protocol", protocol), zap.String("role", role), zap.String("password", password))
		logger.Debug("creating a connection...")

		connectionID, err := v6.CreateConnection(ctx, server, protocol, role, password, datastore)
		if err != nil {
			logger.Error("could not create a connection", zap.Error(err))
			return err
		}

		defer func() {
			logger.Debug("deleting the connection...")

			if err := v6.DeleteConnection(ctx, server, protocol, role, password, datastore, connectionID); err != nil {
				logger.Error("could not delete connection", zap.Error(err))
			}

			logger.Debug("connection deleted!")
		}()

		logger.Debug("connection created", zap.String("connection-id", connectionID))

		logger.Debug("building query...")
		query := fmt.Sprintf("SELECT ?s ?p ?o FROM %s WHERE { ?s ?p ?o }", graph)
		logger.Debug("query built", zap.String("query", query))

		logger.Debug("creating a cursor...")

		cursorID, err := v6.CreateCursor(ctx, server, protocol, role, password, datastore, connectionID, query)
		if err != nil {
			logger.Error("could not create a cursor", zap.Error(err))
			return err
		}

		defer func() {
			logger.Debug("deleting cursor...")

			if err := v6.DeleteCursor(ctx, server, protocol, role, password, datastore, connectionID, cursorID); err != nil {
				logger.Error("could not close the cursor", zap.Error(err))
			}

			logger.Debug("cursor deleted!")
		}()

		logger.Debug("cursor created", zap.String("cursorID", cursorID))

		logger.Debug("opening file for export...")
		f, err := openExportFile(filePath)
		if err != nil {
			logger.Error("could not create export file", zap.Error(err))
			return err
		}

		defer func() {
			logger.Debug("closing file...")

			if err := f.Close(); err != nil {
				logger.Error("could not close file", zap.Error(err))
			}

			logger.Debug("file closed")
		}()

		logger.Info("getting data...")

		var triples map[string]map[string]string
		var tickError error

		readData := func() {
			triples, tickError = v6.ReadWithCursor(ctx, server, protocol, role, password, datastore, connectionID, cursorID, limit)
		}

		doWithTicker(readData, func() {
			logger.Info("still getting data...")
		})

		if tickError != nil {
			logger.Error("could not get data", zap.Error(err))
			return tickError
		}

		logger.Info("got data from the server!")
		logger.Info("writing export data file...")

		writeFile := func() {
			tickError = ttl.Write(triples, f)
		}

		doWithTicker(writeFile, func() {
			logger.Info("still writing file...")
		})

		if tickError != nil {
			return tickError
		}

		logger.Info("export file written!")

		return nil
	}

	return &cmd
}

func openExportFile(path string) (*os.File, error) {
	_, err := os.Stat(path)

	// file exists
	if err == nil {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return nil, err
		}

		if err := file.Truncate(0); err != nil {
			return nil, err
		}

		return file, nil
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return nil, err
	}

	return os.Create(path)
}

func doWithTicker(action func(), onTick func()) {
	tick := time.Tick(time.Second * 1)
	stop := make(chan struct{})
	defer close(stop)

	go func() {
		for {
			select {
			case <-tick:
				onTick()
			case <-stop:
				return
			}
		}
	}()

	action()
}
