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

type (
	cursor interface {
		Read(context.Context) bool
		Close(context.Context) error
	}

	cursorBase struct {
		baseUrl   string
		role      string
		password  string
		operation string
		limit     int
	}

	TripleCursor struct {
		cursorBase
		Data map[string]map[string][]string
		Err  error
	}
)

var _ cursor = &TripleCursor{}

func (t *TripleCursor) Close(ctx context.Context) error {
	logger := utils.LoggerFromContext(ctx).With(zap.String("op", "delete-cursor"), zap.String("url", t.baseUrl))
	client := utils.HttpClientFromContext(ctx)
	logger.Debug("building request...")

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, t.baseUrl, nil)
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return err
	}

	req.Header.Set("Authorization", utils.BasicAuthHeaderValue(t.role, t.password))

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

func (t *TripleCursor) Read(ctx context.Context) bool {
	url := fmt.Sprintf("%s?operation=%s&limit=%d", t.baseUrl, t.operation, t.limit)

	logger := utils.LoggerFromContext(ctx).With(zap.String("op", "advance-cursor"), zap.String("url", url))
	client := utils.HttpClientFromContext(ctx)

	t.Data = make(map[string]map[string][]string)

	logger.Debug("building request...")

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, nil)
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		t.Err = err
		return false
	}

	req.SetBasicAuth(t.role, t.password)
	req.Header.Set("Accept", "text/tab-separated-values")

	logger.Debug("request built", utils.RequestToLoggerFields(req)...)
	logger.Debug("executing request...")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("could not execute request", zap.Error(err))
		t.Err = err
		return false
	}

	defer res.Body.Close()

	logger.Debug("got response from server", zap.String("status", res.Status), zap.Any("header", res.Header))

	if res.StatusCode != http.StatusOK {
		logger.Error("bad response from server", zap.String("status", res.Status))
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error("could not read response body", zap.Error(err))
			bytes = []byte("COULD NOT READ RESPONSE BODY")
		}

		t.Err = fmt.Errorf("bad response from server: %s - %s", res.Status, string(bytes))

		return false
	}

	logger.Debug("processing records...")

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

		if _, ok := t.Data[s]; !ok {
			t.Data[s] = map[string][]string{}
		}

		if _, ok := t.Data[s][p]; !ok {
			t.Data[s][p] = []string{}
		}

		t.Data[s][p] = append(t.Data[s][p], o)

		i++
	}

	res.Body.Close()

	if i < t.limit {
		logger.Debug("no more data to process")
		t.Err = nil
		return false
	}

	logger.Debug("processed triples", zap.Int("count", i))

	t.operation = "advance"

	return true
}

func CreateTripleCursor(ctx context.Context, server, protocol, role, password, datastore, connectionID, query string, limit int) (cursor, error) {
	logger := utils.LoggerFromContext(ctx).With(zap.String("op", "create-cursor"), zap.String("connection-id", connectionID))
	client := utils.HttpClientFromContext(ctx)

	logger.Debug("building url...")

	url := fmt.Sprintf("%s://%s/datastores/%s/connections/%s/cursors", protocol, server, datastore, connectionID)

	logger.Debug("url built", zap.String("url", url))
	logger.Debug("building request...")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(query))
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Authorization", utils.BasicAuthHeaderValue(role, password))
	req.Header.Set("Content-Type", "application/sparql-query")

	logger.Debug("request built", utils.RequestToLoggerFields(req)...)
	logger.Debug("executing request...")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("could not execute request", zap.Error(err))
		return nil, err
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
			return nil, fmt.Errorf("bad response from server: %s - COULD NOT READ RESPONSE BODY: %s", res.Status, err)
		}

		return nil, fmt.Errorf("bad response from server: %s - %s", res.Status, string(bytes))
	}

	logger.Debug("extracting cursor ID from header...")

	cursorID := res.Header.Get("location")
	cursorID = cursorID[strings.LastIndex(cursorID, "/")+1:]

	logger.Info("cursor created", zap.String("cursor", cursorID))

	c := TripleCursor{
		cursorBase: cursorBase{
			baseUrl:   fmt.Sprint(url, "/", cursorID),
			role:      role,
			password:  password,
			operation: "open",
			limit:     limit,
		},
	}

	return &c, nil
}
