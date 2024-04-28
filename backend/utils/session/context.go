package session

import "context"

type key string

var sessionKey key = "session"

func Attach(ctx context.Context, s *session) context.Context {
	return context.WithValue(ctx, sessionKey, s)
}

func Receive(ctx context.Context) *session {
	s := ctx.Value(sessionKey)

	if s == nil {
		panic("session not found in context, did you forget to use the Session middleware?")
	}

	return s.(*session)
}
