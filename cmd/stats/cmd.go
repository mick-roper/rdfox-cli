package stats

import (
	"fmt"

	"net/http"

	"github.com/mick-roper/rdfox-cli/logging"
	"github.com/mick-roper/rdfox-cli/parse"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Cmd() *cobra.Command {
	var cmd cobra.Command
	var datastore string

	cmd.Use = "stats"
	cmd.Short = "get stats for a server or datastore"

	cmd.Flags().StringVar(&datastore, "datastore", "", "The datastore that you want stats for. Leave blank to get server stats.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := logging.GetFromContext(ctx)
		client := utils.HttpClientFromContext(ctx)

		logger.Debug("getting flags...")

		server := cmd.Flags().Lookup("server").Value.String()
		protocol := cmd.Flags().Lookup("protocol").Value.String()
		role := cmd.Flags().Lookup("role").Value.String()
		password := cmd.Flags().Lookup("password").Value.String()

		logger.Debug("got flags", zap.String("server", server), zap.String("protocol", protocol), zap.String("role", role), zap.String("password", password))

		logger.Debug("building url...")

		url := fmt.Sprint(protocol, "://", server)
		if datastore != "" {
			url = fmt.Sprint(url, "/datastores/", datastore)
		}
		url = fmt.Sprint(url, "?component-info=extended")

		logger.Debug("url built", zap.String("url", url))
		logger.Debug("building request")

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logger.Error("could not build request", zap.Error(err))
			return err
		}

		req.Header.Set("Authorization", utils.BasicAuthHeaderValue(role, password))
		req.Header.Set("Accept", "*/*")
		// req.Header.Set("Accept", "application/x.sparql-results+json-abbrev")

		logger.Debug("request built", zap.Stringer("url", req.URL), zap.String("method", req.Method), zap.Any("headers", req.Header))
		logger.Debug("making request...")

		res, err := client.Do(req)
		if err != nil {
			logger.Error("could not make request", zap.Error(err))
			return err
		}

		defer res.Body.Close()

		logger.Debug("got response", zap.String("status", res.Status), zap.Any("headers", res.Header))

		if res.StatusCode != 200 {
			return fmt.Errorf("bad response from server: %s", res.Status)
		}

		stats := parse.Stats(res.Body)

		logger.Info("got stats", zap.Any("stats", stats))

		return nil
	}

	return &cmd
}
