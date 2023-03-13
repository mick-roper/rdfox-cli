package v6

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mick-roper/rdfox-cli/logging"
	"github.com/mick-roper/rdfox-cli/utils"
	"go.uber.org/zap"
)

type statistics map[string]map[string]interface{}

func GetStats(ctx context.Context, server, protocol, role, password, datastore string) (statistics, error) {
	logger := logging.GetFromContext(ctx)
	client := utils.HttpClientFromContext(ctx)

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
		return nil, err
	}

	req.Header.Set("Authorization", utils.BasicAuthHeaderValue(role, password))
	req.Header.Set("Accept", "*/*")
	// req.Header.Set("Accept", "application/x.sparql-results+json-abbrev")

	logger.Debug("request built", zap.Stringer("url", req.URL), zap.String("method", req.Method), zap.Any("headers", req.Header))
	logger.Debug("making request...")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("could not make request", zap.Error(err))
		return nil, err
	}

	defer res.Body.Close()

	logger.Debug("got response", zap.String("status", res.Status), zap.Any("headers", res.Header))

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("bad response from server: %s", res.Status)
	}

	return parseStats(res.Body), nil
}

func parseStats(r io.Reader) statistics {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	stats := statistics{}
	thisComponent := ""
	scanner.Split(bufio.ScanLines)

	i := -1

	for scanner.Scan() {
		i++
		if i == 0 {
			continue
		}

		t := scanner.Text()

		parts := strings.SplitN(t, "\t", 3)
		_, p, v := parts[0], strings.Trim(parts[1], "\""), strings.Trim(parts[2], "\"")

		if p == "Component name" {
			thisComponent = v
			stats[thisComponent] = map[string]interface{}{}
			continue
		}

		stats[thisComponent][p] = v
	}

	return stats
}
