package persistence

import "context"

type key string

var persistenceKey key = "persistence"

var p Persistence = newInMemory()

func Attach(ctx context.Context) context.Context {
	return context.WithValue(ctx, persistenceKey, p)
}

func Receive(ctx context.Context) Persistence {
	p := ctx.Value(persistenceKey)

	if p == nil {
		panic("persistence not found in context, did you forget to attach it?")
	}

	return p.(Persistence)
}
