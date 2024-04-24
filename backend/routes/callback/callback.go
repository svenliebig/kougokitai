package callback

import (
	"net/http"

	"github.com/svenliebig/kougokitai/authenticator"
	"github.com/svenliebig/kougokitai/routes"
	"github.com/svenliebig/kougokitai/utils/session"
)

func init() {
	routes.RegisterRoute("GET /callback", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	auth := authenticator.Use(r.Context())
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	s := session.Save(w, r)

	if s.Get("state") != state {
		http.Error(w, "Invalid state parameter.", http.StatusBadRequest)
		return
	}

	token, err := auth.Exchange(r.Context(), code)

	if err != nil {
		http.Error(w, "Failed to convert an authorization code into a token.", http.StatusUnauthorized)
		return
	}

	idToken, err := auth.VerifyIDToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Failed to verify ID Token.", http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	s.Set("profile", profile)
	s.Set("access_token", token.AccessToken)

	w.Header().Set("Location", "/dashboard")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
