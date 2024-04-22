package rest

import (
	"net/http"
)

type Option interface {
	apply(*http.Request)
}

type WithHeader struct {
	Key   string
	Value string
}

func (o WithHeader) apply(req *http.Request) {
	req.Header.Set(o.Key, o.Value)
}

type WithBearer struct {
	Token string
}

func (o WithBearer) apply(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+o.Token)
}
