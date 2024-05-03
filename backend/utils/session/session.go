package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
)

const sessionCookieName = "id"

var lock = sync.RWMutex{}
var sessions = map[string]*session{}

// TODO: How to invalid a session by user privilege change? I should not have user data in session, without looking into the session data?

// TODO implement a session timeout

type session struct {
	values  map[string]any
	ip      string
	created int64
}

// TODO: this should only be called once to be honest, in the middleware
//
//	      it would be better to expose a function that can be used by the
//	      middleware to get the session, multiple calls to this function
//			   will create multiple sessions
func Save(w http.ResponseWriter, r *http.Request) (s *session) {
	s = get(r)

	if s == nil {
		// TODO: publish event for unique user request
		s = create(r)
	}

	state, err := generateId()

	if err != nil {
		log.Println("error generating session id:", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return nil
	}

	lock.Lock()
	sessions[state] = s
	lock.Unlock()

	setHeader(w, state)

	return
}

func (s *session) Set(key string, value any) {
	s.values[key] = value
}

func Debug() {
	lock.RLock()
	defer lock.RUnlock()

	log.Printf("sessions: %v", len(sessions))
}

func (s *session) Print() {
	keys := make([]string, 0, len(s.values))

	for k := range s.values {
		keys = append(keys, k)
	}

	log.Printf("session values: %v", keys)

	for k, v := range s.values {
		log.Printf("%s: %v", k, v)
	}
}

func (s *session) Get(key string) any {
	value, _ := s.values[key]
	return value
}

func get(r *http.Request) (s *session) {
	cookie, err := r.Cookie(sessionCookieName)

	if errors.Is(err, http.ErrNoCookie) {
		log.Println("no cookie present")
		return
	}

	if err != nil {
		log.Println("error getting cookie:", err)
		return
	}

	lock.RLock()
	s, ok := sessions[cookie.Value]
	lock.RUnlock()

	if !ok {
		log.Printf("no session found for cookie value %q", cookie.Value)
		return
	}

	if s.ip != r.RemoteAddr {
		// TODO: invalidate session because of ip change
		log.Printf("session ip %q does not match request ip %q", s.ip, r.RemoteAddr)
	}

	clear(cookie.Value)

	return s
}

func create(r *http.Request) *session {
	s := session{values: map[string]any{}}
	s.ip = r.RemoteAddr
	s.created = time.Now().Unix()
	return &s
}

func clear(n string) {
	log.Printf("clearing session %q", n)
	lock.Lock()
	delete(sessions, n)
	lock.Unlock()
}

// generates a random session id
func generateId() (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

func setHeader(w http.ResponseWriter, state string) {
	c := &http.Cookie{
		Name:     sessionCookieName,
		Value:    state,
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	w.Header().Add("Set-Cookie", c.String())
	w.Header().Add("Cache-Control", "no-cache=\"Set-Cookie, Set-Cookie2\"")
}
