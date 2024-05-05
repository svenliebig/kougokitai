package callback

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svenliebig/kougokitai/authenticator"
	"github.com/svenliebig/kougokitai/utils/session"
)

func TestCallback(t *testing.T) {
	t.Run("should return bad request when state is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/callback", nil)

		r = r.WithContext(authenticator.Provide(r.Context()))
		r = r.WithContext(session.Attach(r.Context(), session.Save(w, r)))

		handler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Invalid state parameter.\n", w.Body.String())
	})

	t.Run("should fail on token exchange, when the state is random", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/callback?state=hello", nil)

		s := session.Save(w, r)
		s.Set("state", "hello")
		r = r.WithContext(authenticator.Provide(r.Context()))
		r = r.WithContext(session.Attach(r.Context(), s))

		handler(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "Unauthorized request.\n", w.Body.String())
	})
}
