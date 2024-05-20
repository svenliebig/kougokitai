package callback

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/svenliebig/kougokitai/persistence/inmemory"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/stretchr/testify/assert"
	"github.com/svenliebig/kougokitai/authenticator"
	"github.com/svenliebig/kougokitai/authenticator/authenticatortest"
	"github.com/svenliebig/kougokitai/persistence"
	"github.com/svenliebig/kougokitai/utils/session"
	"golang.org/x/oauth2"
)

func new() (*authenticator.Authenticator, *authenticatortest.MockOAuth2Server) {
	server := authenticatortest.NewMockOAuth2Server()
	go server.ListenAndServe()

	provider, err := oidc.NewProvider(
		context.Background(),
		"http://"+server.Addr,
	)

	if err != nil {
		panic("test setup problem: " + err.Error())
	}

	conf := oauth2.Config{
		ClientID:     "client_id",
		ClientSecret: "0wm1GAqlPMuKLaje0sxwNQXhISHP6pqVoG2HgpbKcaZDP1oyBMRDVmKx6zPvNRcO",
		RedirectURL:  "http://localhost:8080/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &authenticator.Authenticator{
		Provider: provider,
		Config:   conf,
	}, server
}

func TestCallbackState(t *testing.T) {
	t.Run("should return bad request when state is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/callback", nil)
		auth, server := new()
		defer server.Close()

		r = r.WithContext(authenticator.Attach(r.Context(), auth))
		r = r.WithContext(session.Attach(r.Context(), session.Save(w, r)))

		handler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Invalid state parameter.\n", w.Body.String())
	})

	t.Run("should return bad request when state does not match the session state", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/callback?state=hello&code=world", nil)
		auth, server := new()
		defer server.Close()

		s := session.Save(w, r)
		s.Set("state", "random")
		r = r.WithContext(authenticator.Attach(r.Context(), auth))
		r = r.WithContext(session.Attach(r.Context(), s))

		handler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Invalid state parameter.\n", w.Body.String())
	})
}

func TestCallbackTokenExchange(t *testing.T) {
	t.Run("should do a token exchange and redirect", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/callback?state=hello&code=world", nil)
		auth, server := new()
		defer server.Close()

		s := session.Save(w, r)
		s.Set("state", "hello")
		r = r.WithContext(authenticator.Attach(r.Context(), auth))
		r = r.WithContext(session.Attach(r.Context(), s))
		r = r.WithContext(persistence.Attach(r.Context()))

		handler(w, r)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		assert.Equal(t, "/welcome", w.Header().Get("Location"))
	})
}
