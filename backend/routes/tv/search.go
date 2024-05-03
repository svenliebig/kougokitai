package tv

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/svenliebig/kougokitai/internal/themoviedb"
	"github.com/svenliebig/kougokitai/routes"
)

var searchUrl = "/tv/search"

func init() {
	routes.RegisterAuthenticatedRoute("GET "+searchUrl, search)
}

func search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	p, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil || p < 1 {
		p = 1
	}

	response, err := themoviedb.SearchTVShows(r.Context(), themoviedb.SearchTVQuery{
		Query: q,
		Page:  p,
	})

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

type PageableResponse[T any] struct {
	Page    int
	Results []T
	Total   int
}
