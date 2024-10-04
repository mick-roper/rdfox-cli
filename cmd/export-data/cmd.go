package exportdata

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	v7 "github.com/mick-roper/rdfox-cli/rdfox/v7"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Cmd() *cobra.Command {
	var (
		cmd       cobra.Command
		datastore string
		filePath  string
		limit     int
		graph     string
		export    string
	)

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

		logger.Debug("building query...")
		var query string

		switch export {
		case "all":
			query = fmt.Sprintf("SELECT ?S ?P ?O FROM <%s> WHERE { ?S ?P ?O }", graph)
		case "explicit":
			query = fmt.Sprintf("SELECT ?S ?P ?O FROM <%s> WHERE { ?S ?P ?O EXPLICIT TRUE }", graph)
		case "implicit":
			query = fmt.Sprintf("SELECT ?S ?P ?O FROM <%s> WHERE { ?S ?P ?O EXPLICIT FALSE }", graph)
		default:
			return errors.New("export must be one of 'all', 'explicit' or 'implicit'")
		}

		logger.Debug("query built", zap.String("query", query))

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

		doneChan := make(chan struct{})
		errChan := make(chan error)

		go func() {
			defer close(doneChan)
			res, err := v7.Query(ctx, protocol, server, role, password, datastore, strings.NewReader(query))
			if err != nil {
				errChan <- err
				return
			}

			defer res.Close()

			logger.Info("writing file...")

			if _, err = io.Copy(f, res); err != nil {
				errChan <- err
				return
			}

			logger.Info("write complete!")
		}()

		go func() {
			t := time.Tick(time.Second * 5)
			for {
				<-t
				logger.Info("writing query response...")
			}
		}()

		select {
		case <-doneChan:
			return nil
		case err := <-errChan:
			return err
		}
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
