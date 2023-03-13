package http

import (
	"context"
	"net/http"

	"github.com/mick-roper/rdfox-cli/utils"
)

var httpClientKey = struct{}{}

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

func ClientFromContext(ctx context.Context) Client {
	if x, ok := utils.GetFromContext(ctx, httpClientKey).(Client); ok {
		return x
	}

	return http.DefaultClient
}

func AddClientToContext(ctx context.Context, client Client) context.Context {
	return utils.AddToContext(ctx, httpClientKey, client)
}
