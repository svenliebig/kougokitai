package dashboard

import (
	"fmt"
	"net/http"

	"github.com/svenliebig/kougokitai/routes"
	"github.com/svenliebig/kougokitai/utils/session"
)

func init() {
	routes.RegisterAuthenticatedRoute("GET /dashboard", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	s := session.Receive(r.Context())
	profile := s.Get("profile")
	w.Write([]byte("Hello, " + fmt.Sprintf("%v", profile)))
}
