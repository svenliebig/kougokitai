package themoviedb

import (
	"context"
)

type Client interface {
	SearchTVShows(context.Context, SearchTVShowsQuery) (s SearchTVShowsResponse, err error)
}

type client struct {
	apiKey string
}

func NewClient(apiKey string) *client {
	return &client{
		apiKey: apiKey,
	}
}
