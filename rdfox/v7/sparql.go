package v7

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mick-roper/rdfox-cli/utils"
)

func Query(ctx context.Context, protocol, server, role, password, datastore string, query io.Reader) (io.ReadCloser, error) {
	logger := utils.LoggerFromContext(ctx)
	client := utils.HttpClientFromContext(ctx)
	url := fmt.Sprintf("%s://%s/datastores/%s/sparql", protocol, server, datastore)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, query)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(role, password)
	req.Header.Set("Content-Type", "application/sparql-query")
	req.Header.Set("Accept", "text/turtle")
	req.Header.Set("Query-Time-Limit", "unlimited")

	logger.Debug("request built", utils.RequestToLoggerFields(req)...)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
