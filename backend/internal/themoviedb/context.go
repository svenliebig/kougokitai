package themoviedb

import (
	"context"
)

type key string

var k string = "themoviedb_client_key"

func Attach(ctx context.Context, client Client) context.Context {
	return context.WithValue(ctx, k, client)
}

func Receive(ctx context.Context) Client {
	return ctx.Value(k).(Client)
}
