package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func clean() {
	sessions = make(map[string]*session)
}

func TestSessionSave(t *testing.T) {
	t.Run("should put a new session into the sessions map with the given ip when no cookie is present", func(t *testing.T) {
		clean()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.RemoteAddr = "helloworld"

		assert.Equal(t, 0, len(sessions))
		Save(w, r)
		assert.Equal(t, 1, len(sessions))

		keys := make([]string, 0, len(sessions))
		for k := range sessions {
			keys = append(keys, k)
		}

		s := sessions[keys[0]]
		assert.Equalf(t, "helloworld", s.ip, "expected ip to be helloworld, got %s", s.ip)
	})

	t.Run("should put a new session into the sessions map with the given ip when a cookie is present but the state is unknown", func(t *testing.T) {
		clean()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.AddCookie(&http.Cookie{
			Name:   sessionCookieName,
			Value:  "random",
			MaxAge: 60,
		})

		r.RemoteAddr = "helloworld"

		assert.Equalf(t, 0, len(sessions), "expected the sessions to be empty before the Save call")
		Save(w, r)
		assert.Equal(t, 1, len(sessions))

		keys := make([]string, 0, len(sessions))
		for k := range sessions {
			keys = append(keys, k)
		}

		s := sessions[keys[0]]
		assert.Equalf(t, "helloworld", s.ip, "expected ip to be helloworld, got %s", s.ip)
	})

	t.Run("should return the session when a cookie is present and request ip is equal to the session ip", func(t *testing.T) {
		clean()

		sessionState := "westeros"
		ip := "helloworld"

		s := &session{
			values: make(map[string]any),
			ip:     ip,
		}
		s.Set("jon", "snow")
		sessions[sessionState] = s

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = ip

		r.AddCookie(&http.Cookie{
			Name:   sessionCookieName,
			Value:  sessionState,
			MaxAge: 60,
		})

		rs := Save(w, r)

		assert.Equalf(t, "snow", rs.Get("jon"), "expected value to be snow, got %s", rs.Get("jon"))
	})

	t.Run("should clear the session when a cookie is present but request ip is not equal to the session ip", func(t *testing.T) {
		clean()

		sessionState := "westeros"
		ip := "helloworld"

		s := &session{
			values: make(map[string]any),
			ip:     ip,
		}
		s.Set("jon", "snow")
		sessions[sessionState] = s

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = "winterfell"

		r.AddCookie(&http.Cookie{
			Name:   sessionCookieName,
			Value:  sessionState,
			MaxAge: 60,
		})

		Save(w, r)

		rs := sessions[sessionState]
		assert.Nilf(t, rs, "expected session to be cleared, but got %v", rs)
	})
}

func TestSessionSaveSetCookie(t *testing.T) {
	t.Run("should set a cookie with the session state as value", func(t *testing.T) {
		clean()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		r.RemoteAddr = "192.1.1.1"

		assert.Equal(t, 0, len(sessions))
		Save(w, r)
		assert.Equal(t, 1, len(sessions))

		cookie := w.Result().Cookies()[0]
		assert.Equalf(t, sessionCookieName, cookie.Name, "expected cookie name to be %s, got %s", sessionCookieName, cookie.Name)
		assert.Equalf(t, sessions[cookie.Value].ip, r.RemoteAddr, "expected ip to be %s, got %s", r.RemoteAddr, sessions[cookie.Value].ip)
		assert.NotEmpty(t, cookie.Value, "expected cookie value to be set")
		assert.True(t, cookie.HttpOnly, "expected cookie to be http only")
		assert.True(t, cookie.Secure, "expected cookie to be secure")
		assert.Equal(t, "/", cookie.Path, "expected cookie path to be /")
		assert.Equal(t, 60*60*24*7, cookie.MaxAge, "expected cookie max age to be 604800")
		assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite, "expected cookie same site to be lax")
	})
}
