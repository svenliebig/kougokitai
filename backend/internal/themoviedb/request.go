package themoviedb

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/svenliebig/kougokitai/internal/rest"
)

const base = "https://api.themoviedb.org/3"

var (
	bearer rest.Option
)

func request(ctx context.Context, path string) (*http.Response, error) {
	apiKey := os.Getenv("THE_MOVIE_DB_API_KEY")
	bearer = rest.WithBearer{
		Token: apiKey,
	}

	res, err := rest.Get(ctx, fmt.Sprintf("%s%s", base, path), bearer, rest.WithHeader{Key: "Accept", Value: "application/json"})

	if err != nil {
		panic(err)
	}

	return res, nil
}
