package middleware

import (
	"net/http"

	"github.com/svenliebig/kougokitai/utils/session"
)

func Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Save(w, r)
		next.ServeHTTP(w, r.WithContext(session.Attach(r.Context(), s)))
	})
}
