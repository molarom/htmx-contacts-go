package htmx

import "context"

type ctxKey int

const hxKey ctxKey = 0

func Get(ctx context.Context) HX {
	v, ok := ctx.Value(hxKey).(HX)
	if !ok {
		return HX{}
	}

	return v
}

func set(ctx context.Context, hx HX) context.Context {
	return context.WithValue(ctx, hxKey, hx)
}
