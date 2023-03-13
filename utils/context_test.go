package utils

import (
	"context"
	"testing"
)

func TestContextUtils(t *testing.T) {
	var key = struct{}{}
	value := "hello, world"

	ctx := context.TODO()

	ctx = AddToContext(ctx, key, value)

	got := GetFromContext(ctx, key)

	if value != got {
		t.Errorf("want = %v, got %v", value, got)
	}
}
