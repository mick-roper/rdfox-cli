package v6

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mick-roper/rdfox-cli/utils"
	"go.uber.org/zap"
)

func CreateConnection(ctx context.Context, server, protocol, role, password, datastore string) (string, error) {
	logger := utils.LoggerFromContext(ctx).With(zap.String("op", "create-connection"))
	client := utils.HttpClientFromContext(ctx)

	logger.Debug("creating url...")

	url := fmt.Sprintf("%s://%s/datastores/%s/connections", protocol, server, datastore)

	logger.Debug("url created", zap.String("url", url))
	logger.Debug("creating request...")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return "", err
	}

	req.Header.Set("Authorization", utils.BasicAuthHeaderValue(role, password))

	logger.Debug("request built", utils.RequestToLoggerFields(req)...)
	logger.Debug("making request...")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("could not make request", zap.Error(err))
		return "", err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Error("could not close response body", zap.Error(err))
		}
	}()

	if res.StatusCode != http.StatusCreated {
		logger.Error("bad response from server", zap.String("status", res.Status))

		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("bad response from server: %s - COULD NOT READ BODY: %s", res.Status, err)
		}

		return "", fmt.Errorf("bad response from server: %s - %s", res.Status, string(bytes))
	}

	c := res.Header.Get("location")
	c = c[strings.LastIndex(c, "/")+1:]

	logger.Info("connection created", zap.String("connection-id", c))

	return c, nil
}

func DeleteConnection(ctx context.Context, server, protocol, role, password, datastore, connectionID string) error {
	logger := utils.LoggerFromContext(ctx).With(zap.String("op", "delete-connection"), zap.String("datastore", datastore), zap.String("connection-id", connectionID))
	client := utils.HttpClientFromContext(ctx)

	logger.Debug("creating url...")

	url := fmt.Sprintf("%s://%s/datastores/%s/connections/%s", protocol, server, datastore, connectionID)

	logger.Debug("url created", zap.String("url", url))
	logger.Debug("creating request...")

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return err
	}

	req.Header.Set("Authorization", utils.BasicAuthHeaderValue(role, password))

	logger.Debug("request built", utils.RequestToLoggerFields(req)...)
	logger.Debug("making request...")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("could not make request", zap.Error(err))
		return err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Error("could not close response body", zap.Error(err))
		}
	}()

	if res.StatusCode != http.StatusNoContent {
		logger.Error("bad response from server", zap.String("status", res.Status))

		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("bad response from server: %s - COULD NOT READ BODY: %s", res.Status, err)
		}

		return fmt.Errorf("bad response from server: %s - %s", res.Status, string(bytes))
	}

	logger.Info("connection closed")

	return nil
}
