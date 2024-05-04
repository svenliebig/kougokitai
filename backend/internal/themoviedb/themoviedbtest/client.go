package themoviedbtest

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/svenliebig/kougokitai/internal/themoviedb"
)

var _ themoviedb.Client = &client{}

type client struct {
	mock.Mock
}

func NewClient() *client {
	return &client{}
}

func (c *client) SearchTVShows(ctx context.Context, query themoviedb.SearchTVShowsQuery) (s themoviedb.SearchTVShowsResponse, err error) {
	args := c.Called(ctx, query)
	return args.Get(0).(themoviedb.SearchTVShowsResponse), args.Error(1)
}
