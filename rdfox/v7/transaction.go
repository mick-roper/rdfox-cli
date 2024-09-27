package v7

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mick-roper/rdfox-cli/utils"
	"go.uber.org/zap"
)

func CreateReadOnlyTransaction(ctx context.Context, server, protocol, role, password, datastore string, connectionID string) error {
	return transactionOp(ctx, server, protocol, role, password, datastore, connectionID, "begin-read-only-transaction")
}

func CreateReadWriteTransaction(ctx context.Context, server, protocol, role, password, datastore string, connectionID string) error {
	return transactionOp(ctx, server, protocol, role, password, datastore, connectionID, "begin-read-write-transaction")
}

func CommitTransaction(ctx context.Context, server, protocol, role, password, datastore string, connectionID string) error {
	return transactionOp(ctx, server, protocol, role, password, datastore, connectionID, "commit-transaction")
}

func RollbackTransaction(ctx context.Context, server, protocol, role, password, datastore string, connectionID string) error {
	return transactionOp(ctx, server, protocol, role, password, datastore, connectionID, "rollback-transaction")
}

func transactionOp(ctx context.Context, server, protocol, role, password, datastore string, connectionID string, transactionOperation string) error {
	logger := utils.LoggerFromContext(ctx).With(zap.String("op", "update transaction"), zap.String("transaction", transactionOperation), zap.String("connection_id", connectionID))
	client := utils.HttpClientFromContext(ctx)

	url := fmt.Sprintf("%s://%s/datastores/%s/connections/%s/transaction?operation=%s", protocol, server, datastore, connectionID, transactionOperation)

	logger.Debug("url created", zap.String("url", url))
	logger.Debug("creating request...")

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, nil)
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return err
	}

	req.SetBasicAuth(role, password)

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

	c := res.Header.Get("location")
	c = c[strings.LastIndex(c, "/")+1:]

	logger.Info("transaction updated", zap.String("connection-id", c))

	return nil
}
