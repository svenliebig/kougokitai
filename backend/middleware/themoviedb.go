package middleware

import (
	"net/http"
	"os"

	"github.com/svenliebig/kougokitai/internal/themoviedb"
)

func ThemoviedbClient(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := themoviedb.NewClient(os.Getenv("THEMOVIEDB_API_KEY"))
		next.ServeHTTP(w, r.WithContext(themoviedb.Attach(r.Context(), c)))
	})
}
