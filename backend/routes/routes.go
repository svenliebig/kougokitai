package routes

import (
	"net/http"
)

var Router = http.NewServeMux()

func RegisterRoute(path string, handler http.HandlerFunc) {
	Router.HandleFunc(path, handler)
}
