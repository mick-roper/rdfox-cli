package v7

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mick-roper/rdfox-cli/utils"
	"go.uber.org/zap"
)

func CreateCursor(ctx context.Context, server, protocol, role, password, datastore, connectionID, query string) (string, error) {
	logger := utils.LoggerFromContext(ctx).With(zap.String("op", "create-cursor"), zap.String("connection-id", connectionID))
	client := utils.HttpClientFromContext(ctx)

	logger.Debug("building url...")

	url := fmt.Sprintf("%s://%s/datastores/%s/connections/%s/cursors", protocol, server, datastore, connectionID)

	logger.Debug("url built", zap.String("url", url))
	logger.Debug("building request...")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(query))
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return "", err
	}

	req.Header.Set("Authorization", utils.BasicAuthHeaderValue(role, password))
	req.Header.Set("Content-Type", "application/sparql-query")

	logger.Debug("request built", utils.RequestToLoggerFields(req)...)
	logger.Debug("executing request...")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("could not execute request", zap.Error(err))
		return "", err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Error("could not close response body", zap.Error(err))
		}
	}()

	logger.Debug("got response from server")

	if res.StatusCode != http.StatusCreated {
		logger.Error("bad response from server", zap.String("status", res.Status))
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error("could not read response body", zap.Error(err))
			return "", fmt.Errorf("bad response from server: %s - COULD NOT READ RESPONSE BODY: %s", res.Status, err)
		}

		return "", fmt.Errorf("bad response from server: %s - %s", res.Status, string(bytes))
	}

	logger.Debug("extracting cursor ID from header...")

	c := res.Header.Get("location")
	c = c[strings.LastIndex(c, "/")+1:]

	logger.Info("cursor created", zap.String("cursor", c))

	return c, nil
}

func DeleteCursor(ctx context.Context, server, protocol, role, password, datastore, connectionID, cursorID string) error {
	logger := utils.LoggerFromContext(ctx).With(zap.String("op", "delete-cursor"), zap.String("connection-id", connectionID), zap.String("cursor-id", cursorID))
	client := utils.HttpClientFromContext(ctx)

	logger.Debug("building url...")

	url := fmt.Sprintf("%s://%s/datastores/%s/connections/%s/cursors/%s", protocol, server, datastore, connectionID, cursorID)

	logger.Debug("url built", zap.String("url", url))
	logger.Debug("building request...")

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return err
	}

	req.Header.Set("Authorization", utils.BasicAuthHeaderValue(role, password))

	logger.Debug("request built", utils.RequestToLoggerFields(req)...)
	logger.Debug("executing request...")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("could not execute request", zap.Error(err))
		return err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Error("could not close response body", zap.Error(err))
		}
	}()

	logger.Debug("got response from server")

	if res.StatusCode != http.StatusNoContent {
		logger.Error("bad response from server", zap.String("status", res.Status))
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error("could not read response body", zap.Error(err))
			return fmt.Errorf("bad response from server: %s - COULD NOT READ RESPONSE BODY: %s", res.Status, err)
		}

		return fmt.Errorf("bad response from server: %s - %s", res.Status, string(bytes))
	}

	logger.Info("cursor closed")

	return nil
}

func ReadWithCursor(ctx context.Context, server, protocol, role, password, datastore, connectionID, cursorID string, limit int, gotData func(map[string]map[string][]string)) error {
	logger := utils.LoggerFromContext(ctx).With(zap.String("op", "advance-cursor"), zap.String("connection-id", connectionID), zap.String("cursor-id", cursorID))
	client := utils.HttpClientFromContext(ctx)

	var read func(op string) error

	read = func(op string) error {
		logger.Debug("building url...")

		url := fmt.Sprintf("%s://%s/datastores/%s/connections/%s/cursors/%s?operation=%s&limit=%d", protocol, server, datastore, connectionID, cursorID, op, limit)

		logger.Debug("url built", zap.String("url", url))
		logger.Debug("building request...")

		req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, nil)
		if err != nil {
			logger.Error("could not build request", zap.Error(err))
			return err
		}

		req.Header.Set("Authorization", utils.BasicAuthHeaderValue(role, password))
		req.Header.Set("Accept", "text/tab-separated-values")

		logger.Debug("request built", utils.RequestToLoggerFields(req)...)
		logger.Debug("executing request...")

		res, err := client.Do(req)
		if err != nil {
			logger.Error("could not execute request", zap.Error(err))
			return err
		}

		defer res.Body.Close()

		logger.Debug("got response from server", zap.String("status", res.Status), zap.Any("header", res.Header))

		if res.StatusCode != http.StatusOK {
			logger.Error("bad response from server", zap.String("status", res.Status))
			bytes, err := io.ReadAll(res.Body)
			if err != nil {
				logger.Error("could not read response body", zap.Error(err))
				return fmt.Errorf("bad response from server: %s - COULD NOT READ RESPONSE BODY: %s", res.Status, err)
			}

			return fmt.Errorf("bad response from server: %s - %s", res.Status, string(bytes))
		}

		logger.Debug("processing records...")

		data := map[string]map[string][]string{}
		var i int
		scanner := bufio.NewScanner(res.Body)
		scanner.Split(bufio.ScanLines)
		scanner.Scan() // always do a dumb scan to skip the first line
		for scanner.Scan() {
			chunks := strings.Split(scanner.Text(), "\t")

			if len(chunks) != 3 {
				logger.Warn("row is the wrong size", zap.Int("row_index", i))
				continue
			}

			s := string(chunks[0])
			p := string(chunks[1])
			o := chunks[2]

			if _, ok := data[s]; !ok {
				data[s] = map[string][]string{}
			}

			if _, ok := data[s][p]; !ok {
				data[s][p] = []string{}
			}

			data[s][p] = append(data[s][p], o)

			i++
		}

		res.Body.Close()

		if i == 0 {
			logger.Debug("no more data to process")
			return nil
		}

		logger.Debug("processed triples", zap.Int("count", i))

		gotData(data)

		logger.Debug("cursor has more data - advancing the cursor")
		return read("advance")
	}

	return read("open")
}
