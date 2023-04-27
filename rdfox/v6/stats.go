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

func GetStats(ctx context.Context, server, protocol, role, password, datastore string) (Statistics, error) {
	logger := utils.LoggerFromContext(ctx)
	client := utils.HttpClientFromContext(ctx)

	logger.Debug("building url...")

	url := fmt.Sprint(protocol, "://", server)
	if datastore != "" {
		url = fmt.Sprint(url, "/datastores/", datastore)
	}
	url = fmt.Sprint(url, "?component-info=extended")

	logger.Debug("url built", zap.String("url", url))
	logger.Debug("building request")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Authorization", utils.BasicAuthHeaderValue(role, password))
	req.Header.Set("Accept", "*/*")

	logger.Debug("request built", utils.RequestToLoggerFields(req)...)
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

func parseStats(r io.Reader) Statistics {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	stats := Statistics{}
	thisComponent := ""
	scanner.Split(bufio.ScanLines)

	i := -1

	for scanner.Scan() {
		i++
		if i == 0 {
			continue
		}

		bytes := scanner.Bytes()
		var chunks [3][]byte
		i := 0

		for _, b := range bytes {
			if b != '\t' {
				chunks[i] = append(chunks[i], b)
				continue
			}

			i++
		}

		p := strings.Trim(string(chunks[1]), "\"")
		v := strings.Trim(string(chunks[2]), "\"")

		if p == "Component name" {
			thisComponent = v
			stats[thisComponent] = map[string]interface{}{}
			continue
		}

		stats[thisComponent][p] = v
	}

	return stats
}
