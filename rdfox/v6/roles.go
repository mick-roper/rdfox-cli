package v6

import (
	"bufio"
	"context"
	"fmt"
	"net/http"

	"github.com/mick-roper/rdfox-cli/utils"
	"go.uber.org/zap"
)

func GetRoles(ctx context.Context, server, protocol, role, password string) ([]string, error) {
	logger := utils.LoggerFromContext(ctx)
	client := utils.HttpClientFromContext(ctx)

	logger.Debug("building url...")

	url := fmt.Sprint(protocol, "://", server, "/roles")

	logger.Debug("url built", zap.String("url", url))
	logger.Debug("building request...")

	req, err := utils.NewRequest(http.MethodGet, url, role, password, nil)
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return nil, err
	}

	logger.Debug("request built", utils.RequestToLoggerFields(req)...)
	logger.Debug("making request...")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("could not make request", zap.Error(err))
		return nil, err
	}

	defer res.Body.Close()

	logger.Debug("got response", zap.String("status", res.Status), zap.Any("headers", res.Header))

	if res.StatusCode != http.StatusOK {
		logger.Error("bad response from server", zap.String("status", res.Status))
		return nil, fmt.Errorf("bad response from server: %s", res.Status)
	}

	logger.Debug("parsing response...")

	roles := []string{}
	scanner := bufio.NewScanner(res.Body)
	scanner.Split(bufio.ScanLines)
	scanner.Scan() // always do this to ignore the first line

	for scanner.Scan() {
		roles = append(roles, scanner.Text())
	}

	logger.Debug("response parsed!")

	return roles, nil
}
