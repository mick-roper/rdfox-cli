package utils

import "context"

type bag map[any]any

var bagKey = struct{}{}

func AddToContext(ctx context.Context, key, value any) context.Context {
	thisBag, ok := ctx.Value(bagKey).(bag)
	if ok {
		thisBag[key] = value
		return ctx
	}

	thisBag = bag{
		key: value,
	}

	return context.WithValue(ctx, bagKey, thisBag)
}

func GetFromContext(ctx context.Context, key any) any {
	thisBag, ok := ctx.Value(bagKey).(bag)
	if ok {
		return thisBag[key]
	}

	return nil
}
