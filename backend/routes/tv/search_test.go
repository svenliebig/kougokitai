package tv

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/svenliebig/kougokitai/internal/themoviedb"
	"github.com/svenliebig/kougokitai/internal/themoviedb/themoviedbtest"
)

func TestRoutesTvSearch(t *testing.T) {
	t.Run("should return a 200 with empty search results", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?q=%s&page=%d", searchUrl, "game", 1), nil)

		mockClient := themoviedbtest.NewClient()
		mockClient.On("SearchTVShows", mock.Anything, themoviedb.SearchTVShowsQuery{
			Page:  1,
			Query: "game",
		}).Return(themoviedb.SearchTVShowsResponse{}, nil)
		r = r.WithContext(themoviedb.Attach(r.Context(), mockClient))

		search(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		mockClient.AssertExpectations(t)
	})

	t.Run("should call the api client with the correct page", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?q=%s&page=%d", searchUrl, "westeros", 5), nil)

		mockClient := themoviedbtest.NewClient()
		mockClient.On("SearchTVShows", mock.Anything, themoviedb.SearchTVShowsQuery{
			Page:  5,
			Query: "westeros",
		}).Return(themoviedb.SearchTVShowsResponse{}, nil)
		r = r.WithContext(themoviedb.Attach(r.Context(), mockClient))

		search(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		mockClient.AssertExpectations(t)
	})

	t.Run("should return a 500 internal server error when the client returns an error", func(t *testing.T) {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?q=%s&page=%d", searchUrl, "westeros", 5), nil)

		mockClient := themoviedbtest.NewClient()
		mockClient.On("SearchTVShows", mock.Anything, themoviedb.SearchTVShowsQuery{
			Page:  5,
			Query: "westeros",
		}).Return(themoviedb.SearchTVShowsResponse{}, fmt.Errorf("error"))
		r = r.WithContext(themoviedb.Attach(r.Context(), mockClient))

		search(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		mockClient.AssertExpectations(t)
	})
}
