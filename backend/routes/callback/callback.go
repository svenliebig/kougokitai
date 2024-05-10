package callback

import (
	"log"
	"net/http"

	"github.com/svenliebig/kougokitai/authenticator"
	"github.com/svenliebig/kougokitai/persistence"
	"github.com/svenliebig/kougokitai/routes"
	"github.com/svenliebig/kougokitai/utils/session"
)

func init() {
	routes.RegisterRoute("GET /callback", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	auth := authenticator.Receive(r.Context())
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	s := session.Receive(r.Context())

	if s.Get("state") != state {
		log.Printf("Invalid state parameter %q expected %q", state, s.Get("state"))
		http.Error(w, "Invalid state parameter.", http.StatusBadRequest)
		return
	}

	token, err := auth.ExchangeCode(r.Context(), code)

	if err != nil {
		http.Error(w, "Unauthorized request.", http.StatusUnauthorized)
		return
	}

	idToken, err := auth.VerifyIDToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Failed to verify ID Token.", http.StatusInternalServerError)
		return
	}

	var profile authenticator.Profile
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	s.Set("profile", profile)
	s.Set("access_token", token.AccessToken)

	p := persistence.Receive(r.Context())

	exists, err := p.UserExists(profile.Id)

	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	if !exists {
		p.CreateUser(profile.Id)
		w.Header().Set("Location", "/welcome")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Location", "/dashboard")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
