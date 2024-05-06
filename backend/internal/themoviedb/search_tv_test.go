package themoviedb

import (
	"testing"
)

func TestSearchTVQuery(t *testing.T) {
	t.Run("should transform the query", func(t *testing.T) {
		query := SearchTVShowsQuery{
			Query: "The Simpsons",
		}

		got := query.String()
		want := "?include_adult=false&page=1&query=The+Simpsons"

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
