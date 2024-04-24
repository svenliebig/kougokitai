package middleware

import "net/http"

func Combine(handlers ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		for i := len(handlers) - 1; i >= 0; i-- {
			next = handlers[i](next)
		}
		return next
	}
}
