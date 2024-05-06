package themoviedb

import (
	"context"
	"fmt"
	"net/http"

	"github.com/svenliebig/kougokitai/utils/rest"
)

const base = "https://api.themoviedb.org/3"

var (
	bearer rest.Option
)

func (c *client) request(ctx context.Context, path string) (*http.Response, error) {
	apiKey := c.apiKey
	bearer = rest.WithBearer{
		Token: apiKey,
	}

	res, err := rest.Get(ctx, fmt.Sprintf("%s%s", base, path), bearer, rest.WithHeader{Key: "Accept", Value: "application/json"})

	if err != nil {
		panic(err)
	}

	return res, nil
}
