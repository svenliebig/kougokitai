package authenticator

import (
	"context"
	"errors"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

type AuthenticatorI interface {
	ExchangeCode(context.Context, string) (*oauth2.Token, error)
	VerifyIDToken(context.Context, *oauth2.Token) (*oidc.IDToken, error)
}

type Profile struct {
	Id        string `json:"sid"`
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	Audience  string `json:"aud"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	Name      string `json:"name"`
	Nichname  string `json:"nickname"`
	Picture   string `json:"picture"`
}

type authenticatorKey string

var key authenticatorKey = "authenticator"

func Attach(ctx context.Context, auth AuthenticatorI) context.Context {
	return context.WithValue(ctx, key, auth)
}

func Receive(ctx context.Context) AuthenticatorI {
	return ctx.Value(key).(AuthenticatorI)
}

// New instantiates the *Authenticator.
func New() (AuthenticatorI, error) {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+os.Getenv("AUTH0_DOMAIN")+"/",
	)

	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

func (a *Authenticator) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return a.Exchange(ctx, code)
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
