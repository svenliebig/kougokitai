package persistence

import "context"

type key string

var persistenceKey key = "persistence"

func Attach(ctx context.Context) context.Context {
	return context.WithValue(ctx, persistenceKey, newInMemory())
}

func Receive(ctx context.Context) Persistence {
	p := ctx.Value(persistenceKey)

	if p == nil {
		panic("persistence not found in context, did you forget to attach it?")
	}

	return p.(Persistence)
}
