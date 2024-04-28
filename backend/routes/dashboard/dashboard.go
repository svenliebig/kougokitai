package dashboard

import (
	"fmt"
	"net/http"

	"github.com/svenliebig/kougokitai/routes"
	"github.com/svenliebig/kougokitai/utils/session"
)

func init() {
	routes.RegisterRoute("GET /dashboard", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	s := session.Receive(r.Context())
	profile := s.Get("profile")

	if profile == nil {
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("Hello, " + fmt.Sprintf("%v", profile)))
}
