package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

var httpClientKey = struct{}{}

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

func HttpClientFromContext(ctx context.Context) Client {
	if x, ok := getFromContext(ctx, httpClientKey).(Client); ok {
		return x
	}

	return http.DefaultClient
}

func AddHttpClientToContext(ctx context.Context, client Client) context.Context {
	return addToContext(ctx, httpClientKey, client)
}

func BasicAuthHeaderValue(username, password string) string {
	plaintext := fmt.Sprint(username, ":", password)
	encoded := string(base64.StdEncoding.EncodeToString([]byte(plaintext)))
	return fmt.Sprint("Basic ", encoded)
}

func RequestToLoggerFields(req *http.Request) []zap.Field {
	return []zap.Field{
		zap.Stringer("url", req.URL),
		zap.String("method", req.Method),
		zap.Any("headers", req.Header),
	}
}

func ResponseToLoggerFields(res *http.Response) []zap.Field {
	return []zap.Field{
		zap.String("status", res.Status),
		zap.Any("headers", res.Header),
	}
}

func NewRequest(method, url, role, password string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", BasicAuthHeaderValue(role, password))

	return req, nil
}
