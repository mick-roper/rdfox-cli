package v6

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mick-roper/rdfox-cli/utils"
	"go.uber.org/zap"
)

func ImportAxioms(ctx context.Context, protocol, server, role, password, datastore, srcGraph, dstGraph string) error {
	logger := utils.LoggerFromContext(ctx)
	client := utils.HttpClientFromContext(ctx)

	logger.Debug("building url...")

	url := fmt.Sprintf("%s://%s/datastores/%s/content?operation=add-axioms&source-graph=%s&destination-graph=%s", protocol, server, datastore, srcGraph, dstGraph)

	logger.Debug("url built", zap.String("url", url))
	logger.Debug("building request...")

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, nil)
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return err
	}

	req.Header.Set("Authorization", utils.BasicAuthHeaderValue(role, password))

	logger.Debug("request built", utils.RequestToLoggerFields(req)...)

	logger.Debug("making request...")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("could not get response", zap.Error(err))
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		payload, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("bad response from server: %s - COULD NOT READ RESPONSE: %s", res.Status, err)
		}
		return fmt.Errorf("bad response from server: %s %s", res.Status, string(payload))
	}

	scanner := bufio.NewScanner(res.Body)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		logger.Info("success", zap.String("data", s))
	}

	return nil
}
