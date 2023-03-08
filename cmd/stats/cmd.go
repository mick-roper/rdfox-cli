package stats

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mick-roper/rdfox-cli/logging"
	"github.com/mick-roper/rdfox-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Cmd() *cobra.Command {
	var datastore string

	cmd := cobra.Command{
		Use:   "stats",
		Short: "prints statistics",
		Long:  "prints statistics about your RDFox server or datastores",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logging.GetFromContext(ctx)
			role := cmd.Flag("role").Value.String()
			password := cmd.Flag("password").Value.String()
			protocol := cmd.Flag("protocol").Value.String()
			server := cmd.Flag("server").Value.String()

			url := fmt.Sprint(protocol, "://", server)

			if datastore != "" {
				url = fmt.Sprint(url, "/datastores/", datastore)
			}

			url = fmt.Sprint(url, "?component-info=extended")
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				logger.Error("could not build request", zap.Error(err))
				return err
			}

			req.Header.Set("Authorization", utils.ToBasicAuth(role, password))

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				logger.Error("request failed", zap.Error(err))
				return err
			}

			defer res.Body.Close()

			if res.StatusCode != 200 {
				logger.Error("bad response from server", zap.String("status-code", res.Status))
				return errors.New("bad status")
			}

			b, err := io.ReadAll(res.Body)
			if err != nil {
				logger.Error("could not read response", zap.Error(err))
				return err
			}

			logger.Info("got response", zap.ByteString("data", b))

			return nil
		},
	}

	cmd.Flags().StringVar(&datastore, "datastore", "", "the datastore to get statistics about")

	return &cmd
}
