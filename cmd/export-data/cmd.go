package exportdata

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	var export string

	cmd.Use = "export-data"
	cmd.Short = "export data from the database"
	cmd.Long = "TODO: write sometihn inspiring here!"

	cmd.Flags().StringVar(&datastore, "datastore", "", "the datastore that contains the data you want to export")
	cmd.Flags().StringVar(&filePath, "file", "export.ttl", "the file that the exported data will be written to")
	cmd.Flags().IntVar(&limit, "limit", 5000, "the maximum number of triples to return in a single cursor request")
	cmd.Flags().StringVar(&graph, "graph", "", "the graph that contains the data you want to export")
	cmd.Flags().StringVar(&export, "export", "all", "the types of facts to export: options are 'all', 'explicit' or 'implicit'")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if datastore == "" {
			return errors.New("datastore is unset")
		}

		if graph == "" {
			return errors.New("graph is unset")
		}

		if !(export == "all" || export == "explicit" || export == "implicit") {
			return errors.New("export must be one of 'all', 'explicit' or 'implicit'")
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
		var query string

		switch export {
		case "all":
			query = fmt.Sprintf("SELECT ?s ?p ?o FROM <%s> WHERE { ?s ?p ?o }", graph)
		case "explicit":
			query = fmt.Sprintf("SELECT ?s ?p ?o FROM <%s> WHERE { ?s ?p ?o EXPLICIT TRUE }", graph)
		case "implicit":
			query = fmt.Sprintf("SELECT ?s ?p ?o FROM <%s> WHERE { ?s ?p ?o EXPLICIT FALSE }", graph)
		default:
			return errors.New("unsupported export option")
		}

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

		dataChan := make(chan map[string]map[string][]string)
		readDoneChan := make(chan struct{})
		writeDoneChan := make(chan struct{})

		write := func() {
			defer close(writeDoneChan)
			for {
				select {
				case triples := <-dataChan:
					writeFile := func() error {
						if err := ttl.Write(triples, f); err != nil {
							return err
						}

						return nil
					}

					logger.Info("writing data to file...")

					if err := utils.DoWithTicker(writeFile, func() {
						logger.Info("still writing file...")
					}); err != nil {
						logger.Error("could not write data", zap.Error(err))
						return
					}

					logger.Info("write complete")
				case <-readDoneChan:
					return
				}
			}
		}

		go write()

		readData := func() error {
			defer close(readDoneChan)
			defer close(dataChan)

			handle := func(data map[string]map[string][]string) {
				dataChan <- data
			}

			logger.Info("reading data from the server...")

			if err := v6.ReadWithCursor(ctx, server, protocol, role, password, datastore, connectionID, cursorID, limit, handle); err != nil {
				return err
			}

			logger.Info("read complete")

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
