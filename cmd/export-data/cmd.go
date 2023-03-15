package exportdata

import (
	"errors"
	"fmt"
	"os"

	v6 "github.com/mick-roper/rdfox-cli/rdfox/v6"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Cmd() *cobra.Command {
	var cmd cobra.Command

	var datastore string
	var filePath string
	var limit int
	var query string

	cmd.Use = "export-data"
	cmd.Short = "export data from the database"
	cmd.Long = "TODO: write sometihn inspiring here!"

	cmd.Flags().StringVar(&datastore, "datastore", "", "the datastore that contains the data you want to export")
	cmd.Flags().StringVar(&filePath, "file", "export.ttl", "the file that the exported data will be written to")
	cmd.Flags().IntVar(&limit, "limit", 5000, "the maximum number of triples to return in a single cursor request")
	cmd.Flags().StringVar(&query, "query", "", "the sparql query that will be used to select data to export")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if datastore == "" {
			return errors.New("datastore is unset")
		}

		if query == "" {
			return errors.New("query is unset")
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

		logger.Debug("getting data...")

		ttl, err := v6.ReadWithCursor(ctx, server, protocol, role, password, datastore, connectionID, cursorID, limit)
		if err != nil {
			logger.Error("could not read data", zap.Error(err))
			return err
		}

		logger.Debug("got data from the server!")
		logger.Debug("writing file...")

		for s, duples := range ttl {
			logger.Debug("writing subject to file...", zap.String("subject", s))

			if _, err := f.WriteString(s); err != nil {
				return err
			}

			var x int

			for p, o := range duples {
				x++

				str := fmt.Sprint("\n\t", p, "\t", o)
				if _, err := f.WriteString(str); err != nil {
					return err
				}

				if x < len(duples) {
					f.WriteString(";")
				} else {
					f.WriteString(".")
				}
			}

			f.WriteString("\n")
		}

		logger.Debug("file written!")

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

	return os.Create(path)
}
