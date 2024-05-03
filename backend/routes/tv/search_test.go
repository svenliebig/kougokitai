package tv

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutesTvSearch(t *testing.T) {
	t.Run("should return a 200 with the search results", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?q=%s&p=%d", searchUrl, "game", 1), nil)
		w := httptest.NewRecorder()

		search(w, r)

		// can I mock themoviedb.SearchTVShows?

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})
}
