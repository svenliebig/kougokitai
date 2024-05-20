package authenticatortest

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ProviderJSON struct {
	Issuer        string   `json:"issuer"`
	AuthURL       string   `json:"authorization_endpoint"`
	TokenURL      string   `json:"token_endpoint"`
	DeviceAuthURL string   `json:"device_authorization_endpoint"`
	JWKSURL       string   `json:"jwks_uri"`
	UserInfoURL   string   `json:"userinfo_endpoint"`
	Algorithms    []string `json:"id_token_signing_alg_values_supported"`
}

// tokenJSON is the struct representing the HTTP response from OAuth2
// providers returning a token or error in JSON form.
// https://datatracker.ietf.org/doc/html/rfc6749#section-5.1
type TokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"` // at least PayPal returns string, while most return number
	IdToken      string `json:"id_token"`   // optional
	// error fields
	// https://datatracker.ietf.org/doc/html/rfc6749#section-5.2
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
}

type JWK struct {
	Alg string   `json:"alg"`
	Kty string   `json:"kty"`
	Use string   `json:"use"`
	Kid string   `json:"kid"`
	X5t string   `json:"x5t"`
	X5c []string `json:"x5c"`
	N   string   `json:"n"`
	E   string   `json:"e"`
}

type JWKSJSON struct {
	Keys []JWK `json:"keys"`
}

type Route string

const (
	OpenIDConfiguration Route = "/.well-known/openid-configuration"
	Token               Route = "/token"
	Jwks                Route = "/jwks"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// creates a mock server for testing services and functions that use OAuth internally.
//
// by default the server will return valid responses for all requests, it's possible
// to override the default behavior by providing a custom handler for a specific path.
type MockOAuth2Server struct {
	http.Server

	responses map[string]HandlerFunc
	mutex     sync.Mutex
}

func NewMockOAuth2Server() *MockOAuth2Server {
	server := &MockOAuth2Server{
		Server: http.Server{
			Addr: "localhost:9123",
		},
		responses: make(map[string]HandlerFunc, 0),
	}

	mux := http.NewServeMux()

	mux.HandleFunc(string(OpenIDConfiguration), func(w http.ResponseWriter, r *http.Request) {
		server.handleOpenIDConfiguration(w, r)
	})

	mux.HandleFunc(string(Token), func(w http.ResponseWriter, r *http.Request) {
		server.handleToken(w, r)
	})

	mux.HandleFunc(string(Jwks), func(w http.ResponseWriter, r *http.Request) {
		server.handleJwks(w, r)
	})

	server.Handler = mux

	return server
}

func (s *MockOAuth2Server) handleOpenIDConfiguration(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	response, ok := s.responses[r.URL.Path]

	if !ok {
		w.Header().Set("Content-Type", "application/json")
		c := ProviderJSON{
			Issuer:        "http://" + s.Addr,
			AuthURL:       "http://" + s.Addr + "/auth",
			TokenURL:      "http://" + s.Addr + string(Token),
			DeviceAuthURL: "http://" + s.Addr + "/device",
			JWKSURL:       "http://" + s.Addr + string(Jwks),
			UserInfoURL:   "http://" + s.Addr + "/userinfo",
			Algorithms:    []string{"RS256"},
		}
		json.NewEncoder(w).Encode(c)
		return
	}

	response(w, r)
}

func (s *MockOAuth2Server) handleToken(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	response, ok := s.responses[r.URL.Path]

	if !ok {
		w.Header().Set("Content-Type", "application/json")

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"iss": "http://" + s.Addr,
			"aud": "client_id",
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		_, filename, _, _ := runtime.Caller(0)
		keyData, err := os.ReadFile(filepath.Join(filepath.Dir(filename), "sample_key"))

		if err != nil {
			panic("err reading key: " + err.Error())
		}

		rsa, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)

		if err != nil {
			panic("err parsing key: " + err.Error())
		}

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(rsa)

		if err != nil {
			panic("err creating and sigining token: " + err.Error())
		}

		c := TokenJSON{
			AccessToken:  "access_token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
			RefreshToken: "refresh_token",
			IdToken:      tokenString,
		}
		json.NewEncoder(w).Encode(c)
		return
	}

	response(w, r)
}

func (s *MockOAuth2Server) handleJwks(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	response, ok := s.responses[r.URL.Path]

	if !ok {
		w.Header().Set("Content-Type", "application/json")

		_, filename, _, _ := runtime.Caller(0)
		keyData, err := os.ReadFile(filepath.Join(filepath.Dir(filename), "sample_key.pub"))

		if err != nil {
			panic("err reading key: " + err.Error())
		}

		rsa, err := jwt.ParseRSAPublicKeyFromPEM(keyData)

		if err != nil {
			panic("err parsing key: " + err.Error())
		}

		n := base64.RawURLEncoding.EncodeToString(rsa.N.Bytes())
		e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(rsa.E)).Bytes())

		c := JWKSJSON{
			Keys: []JWK{
				{
					Alg: "RS256",
					Kty: "RSA",
					Use: "sig",
					Kid: "key_id",
					N:   n,
					E:   e,
				},
			},
		}

		json.NewEncoder(w).Encode(c)
		return
	}

	response(w, r)
}
