package callback

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/svenliebig/kougokitai/authenticator"
	"github.com/svenliebig/kougokitai/utils/session"
	"golang.org/x/oauth2"
)

type MockAuthenticator struct {
	authenticator.Authenticator
	mock.Mock
}

func (m *MockAuthenticator) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	args := m.Called(code)
	r1 := args.Get(0)

	if r1 == nil {
		return nil, args.Error(1)
	}

	return r1.(*oauth2.Token), args.Error(1)
}

func (m *MockAuthenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	args := m.Called(token)
	return args.Get(0).(*oidc.IDToken), args.Error(1)
}

func TestCallback(t *testing.T) {
	t.Run("should return bad request when state is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/callback", nil)

		r = r.WithContext(authenticator.Attach(r.Context(), &MockAuthenticator{}))
		r = r.WithContext(session.Attach(r.Context(), session.Save(w, r)))

		handler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Invalid state parameter.\n", w.Body.String())
	})

	t.Run("should fail on token exchange, when the state is random", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/callback?state=hello&code=world", nil)

		s := session.Save(w, r)
		s.Set("state", "hello")
		mock := &MockAuthenticator{}
		r = r.WithContext(authenticator.Attach(r.Context(), mock))
		r = r.WithContext(session.Attach(r.Context(), s))

		mock.On("ExchangeCode", "world").Return(nil, errors.New("unauthorized"))

		handler(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "Unauthorized request.\n", w.Body.String())
	})
}
