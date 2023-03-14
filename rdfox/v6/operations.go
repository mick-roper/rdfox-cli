package v6

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mick-roper/rdfox-cli/utils"
	"go.uber.org/zap"
)

func ImportAxioms(ctx context.Context, protocol, server, role, password, datastore, srcGraph, dstGraph string) error {
	logger := utils.LoggerFromContext(ctx)
	client := utils.HttpClientFromContext(ctx)

	logger.Debug("building url...")

	url := fmt.Sprint(protocol, "://", server, "/datastores/", datastore, "/content?operation=add-axioms&source-graph=", srcGraph, "&detination-graph=", dstGraph)

	logger.Debug("url built", zap.String("url", url))
	logger.Debug("building request...")

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, nil)
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return err
	}

	req.Header.Set("Authorization", utils.BasicAuthHeaderValue(role, password))

	logger.Debug("request built", zap.Stringer("url", req.URL), zap.Any("headers", req.Header), zap.String("method", req.Method))

	logger.Debug("making request...")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("could not get response", zap.Error(err))
		return err
	}

	defer res.Body.Close()

	logger.Debug("got response", zap.String("status", res.Status))

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from server: %s", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("could not read response", zap.Error(err))
		return err
	}

	logger.Info("response from server", zap.ByteString("data", bytes))

	return nil
}
