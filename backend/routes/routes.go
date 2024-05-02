package routes

import (
	"net/http"

	"github.com/svenliebig/kougokitai/utils/session"
)

var Router = http.NewServeMux()

func RegisterRoute(path string, handler http.HandlerFunc) {
	Router.HandleFunc(path, handler)
}

func RegisterAuthenticatedRoute(path string, handler http.HandlerFunc) {
	Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		s := session.Receive(r.Context())
		profile := s.Get("profile")

		if profile == nil {
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		handler(w, r)
	})
}
