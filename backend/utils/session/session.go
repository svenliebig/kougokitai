package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
)

const sessionCookieName = "session"

var sessions = map[string]session{}

type session struct {
	values map[string]any
}

// TODO error handling
func Save(w http.ResponseWriter, r *http.Request) *session {
	cookie, err := r.Cookie(sessionCookieName)

	if err == nil && cookie != nil {
		return getOrCreate(cookie.Value)
	}

	if !errors.Is(err, http.ErrNoCookie) {
		log.Println("Error getting cookie:", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return nil
	}

	state, err := generateRandomState()

	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return nil
	}

	w.Header().Set("Set-Cookie", "session="+state+"; HttpOnly; Secure")
	return getOrCreate(state)
}

func (s *session) Set(key string, value any) {
	s.values[key] = value
}

func (s *session) Get(key string) (any) {
	value, _ := s.values[key]
	return value
}

func getOrCreate(name string) *session {
	s, ok := sessions[name]

	if !ok {
		s = session{values: map[string]any{}}
		sessions[name] = s
	}

	return &s
}

func generateRandomState() (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
