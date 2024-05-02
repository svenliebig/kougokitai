package middleware

import (
	"net/http"

	"github.com/svenliebig/kougokitai/persistence"
)

func Persistence(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(persistence.Attach(r.Context())))
	})
}
