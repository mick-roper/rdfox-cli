package compact

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Cmd() *cobra.Command {
	var cmd cobra.Command
	var datastore string

	cmd.Use = "compact"
	cmd.Short = "compacts the database"
	cmd.Long = "todo: write a long description"

	cmd.Flags().StringVar(&datastore, "datastore", "", "the datastore to compact")

	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		logger := utils.LoggerFromContext(ctx)
		client := utils.HttpClientFromContext(ctx)

		rootFlags := utils.RootCommandFlags(cmd)

		logger.Debug("building url...")

		url := fmt.Sprintf("%s://%s/commands", rootFlags.Protocol, rootFlags.Server)

		logger.Debug("url built", zap.String("url", url))
		logger.Debug("building command...")

		command := fmt.Sprintf("active %s\ncompact", datastore)

		logger.Debug("command built", zap.String("command", command))
		logger.Debug("building request")

		req, err := utils.NewRequest(http.MethodPost, url, rootFlags.Role, rootFlags.Password, strings.NewReader(command))
		if err != nil {
			logger.Error("could not build request", zap.Error(err))
			return err
		}

		req.Header.Set("Content-Type", "text/plain")

		logger.Debug("request built!", utils.RequestToLoggerFields(req)...)
		logger.Debug("making request...")

		res, err := client.Do(req)
		if err != nil {
			logger.Error("request failed", zap.Error(err))
			return err
		}

		defer res.Body.Close()
		logger.Debug("request complete!", utils.ResponseToLoggerFields(res)...)

		payload, _ := io.ReadAll(res.Body)

		if res.StatusCode > 299 {
			logger.Error("bad status from server", zap.ByteString("response", payload))
			return fmt.Errorf("bad response from server: %s", res.Status)
		}

		logger.Info("response from server", zap.ByteString("data", payload))

		return nil
	}

	return &cmd
}
