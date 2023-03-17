package v6

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

	logger.Debug("got response", utils.ResponseToLoggerFields(res)...)

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

func CreateRole(ctx context.Context, server, protocol, role, password, newRoleName, newRolePassword string) error {
	logger := utils.LoggerFromContext(ctx)
	client := utils.HttpClientFromContext(ctx)

	logger.Debug("building url...")

	url := fmt.Sprint(protocol, "://", server, "/roles/", newRoleName)

	logger.Debug("url built", zap.String("url", url))
	logger.Debug("building request...")

	req, err := utils.NewRequest(http.MethodPost, url, role, password, strings.NewReader(newRolePassword))
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return err
	}

	req.Header.Set("Content-Type", "text/plain")

	logger.Debug("request built", utils.RequestToLoggerFields(req)...)
	logger.Debug("making request...")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("could not make request", zap.Error(err))
		return err
	}

	defer res.Body.Close()

	logger.Debug("got response", utils.ResponseToLoggerFields(res)...)

	if res.StatusCode != http.StatusCreated {
		logger.Error("bad response from server", zap.String("status", res.Status))
		return fmt.Errorf("bad response from server: %s", res.Status)
	}

	return nil
}

func DeleteRole(ctx context.Context, server, protocol, role, password, roleToDelete string) error {
	logger := utils.LoggerFromContext(ctx)
	client := utils.HttpClientFromContext(ctx)

	logger.Debug("building url...")

	url := fmt.Sprint(protocol, "://", server, "/roles/", roleToDelete)

	logger.Debug("url built", zap.String("url", url))
	logger.Debug("building request...")

	req, err := utils.NewRequest(http.MethodDelete, url, role, password, nil)
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return err
	}

	logger.Debug("request built", utils.RequestToLoggerFields(req)...)
	logger.Debug("making request...")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("could not make request", zap.Error(err))
		return err
	}

	defer res.Body.Close()

	logger.Debug("got response", utils.ResponseToLoggerFields(res)...)

	if res.StatusCode != http.StatusNoContent {
		logger.Error("bad response from server", zap.String("status", res.Status))
		return fmt.Errorf("bad response from server: %s", res.Status)
	}

	return nil
}

func GrantDatastorePrivileges(ctx context.Context, server, protocol, role, password, targetRole, accessTypes string) error {
	return updateDatastoreAccessTypes(ctx, server, protocol, role, password, targetRole, "add", accessTypes)
}

func RevokeDatastorePrivileges(ctx context.Context, server, protocol, role, password, targetRole, accessTypes string) error {
	return updateDatastoreAccessTypes(ctx, server, protocol, role, password, targetRole, "delete", accessTypes)
}

func updateDatastoreAccessTypes(ctx context.Context, server, protocol, role, password, targetRole, operation, accessTypes string) error {
	logger := utils.LoggerFromContext(ctx)
	client := utils.HttpClientFromContext(ctx)

	if operation != "add" && operation != "delete" {
		return errors.New("only 'add' and 'delete' operations are supported")
	}

	logger.Debug("building url...")

	url := fmt.Sprint(protocol, "://", server, "/roles/", targetRole, "/privileges?operation=", operation)

	logger.Debug("url built", zap.String("url", url))
	logger.Debug("building request...")

	bodyString := fmt.Sprint("resource-specifier=|datastores&access-types=", accessTypes)

	req, err := utils.NewRequest(http.MethodPatch, url, role, password, strings.NewReader(bodyString))
	if err != nil {
		logger.Error("could not build request", zap.Error(err))
		return err
	}

	req.Header.Set("Content-Type", "x-www-form-urlencoded")

	logger.Debug("request built", utils.RequestToLoggerFields(req)...)
	logger.Debug("making request...")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("could not make request", zap.Error(err))
		return err
	}

	defer res.Body.Close()

	logger.Debug("got response", utils.ResponseToLoggerFields(res)...)

	if res.StatusCode != http.StatusOK {
		logger.Error("bad response from server", zap.String("status", res.Status))
		return fmt.Errorf("bad response from server: %s", res.Status)
	}

	return nil
}