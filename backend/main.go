package main

import (
	"context"
	"fmt"
	"os"

	"github.com/svenliebig/env"
)

func main() {
	err := env.Load()

	if err != nil {
		panic(err)
	}

	fmt.Println(os.Getenv("THE_MOVIE_DB_API_KEY"))

	s := Server[User, Principal]{
		Authentication: auth{},
		Authorization:  auth{},
	}

	r := Route{
		Method:                 "GET",
		Path:                   "/movies",
		AuthenticationRequired: true,
	}

	s.Register(r, r)

	s.Start(context.Background())
}

type Server[A any, P any] struct {
	Authentication Authentication[A]
	Authorization  Authorization[A, P]
}

func (s Server[A, P]) Start(ctx context.Context) error {
	if s.Authentication != nil {
		ctx = s.Authentication.WithAuthentication(ctx)
	}

	return nil
}

func (s Server[A, P]) Register(r ...Route) error {
	return nil
}

type Request struct{}

type Authentication[T any] interface {
	IsAuthenticated(context.Context, Request) bool
	WithAuthentication(context.Context) context.Context
	GetAuthentication(context.Context) T
}

type Authorization[T any, P any] interface {
	GetPermissions(context.Context, T) P
}

type User struct {
	name string
	jwt  string
}

type Principal struct {
	role        string
	permissions []string
}

var _ Authentication[User] = auth{}
var _ Authorization[User, Principal] = auth{}

type auth struct {
}

// I don't like this anymore, this should be a package.

func (a auth) IsAuthenticated(ctx context.Context, r Request) bool {
	return true
}

func (a auth) WithAuthentication(ctx context.Context) context.Context {
	return ctx
}

func (a auth) GetAuthentication(ctx context.Context) User {
	return User{}
}

func (a auth) GetPermissions(ctx context.Context, user User) Principal {
	return Principal{}
}

type Route struct {
	Method                 string
	Path                   string
	AuthenticationRequired bool
}

type Predicate interface{}

type Middleware interface{}
