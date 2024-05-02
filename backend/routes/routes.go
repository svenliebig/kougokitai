package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/svenliebig/kougokitai/authenticator"
	"github.com/svenliebig/kougokitai/utils/session"
)

var Router = http.NewServeMux()

func RegisterRoute(path string, handler http.HandlerFunc) {
	Router.HandleFunc(path, handler)
}

func RegisterAuthenticatedRoute(path string, handler http.HandlerFunc) {
	Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		s := session.Receive(r.Context())
		// TODO I don't put a pointer into the map, but it seems like I get a pointer out
		profile := s.Get("profile")

		if profile == nil {
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		p := profile.(authenticator.Profile)

		if p.ExpiresAt < time.Now().Unix() {
			// TODO handle different, we should invalidate the session
			// TODO is this the right place to handle this? I would see it in the session, but the session is stupid and does not know what data it holds
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		handler(w, r)
	})
}
