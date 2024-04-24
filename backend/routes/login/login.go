package login

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/svenliebig/kougokitai/authenticator"
	"github.com/svenliebig/kougokitai/routes"
	"github.com/svenliebig/kougokitai/utils/session"
)

func init() {
	routes.RegisterRoute("GET /login", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomState()

	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	s := session.Save(w, r)
	s.Set("state", state)

	auth := authenticator.Use(r.Context())

	w.Header().Set("Location", auth.AuthCodeURL(state))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
