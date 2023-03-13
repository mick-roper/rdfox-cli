package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
)

var httpClientKey = struct{}{}

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

func HttpClientFromContext(ctx context.Context) Client {
	if x, ok := GetFromContext(ctx, httpClientKey).(Client); ok {
		return x
	}

	return http.DefaultClient
}

func AddClientToContext(ctx context.Context, client Client) context.Context {
	return AddToContext(ctx, httpClientKey, client)
}

func BasicAuthHeaderValue(username, password string) string {
	plaintext := fmt.Sprint(username, ":", password)
	encoded := string(base64.StdEncoding.EncodeToString([]byte(plaintext)))
	return fmt.Sprint("Basic ", encoded)
}
