package themoviedb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

type response struct {
	Page         int      `json:"page"`
	Results      []Series `json:"results"`
	TotalPages   int      `json:"total_pages"`
	TotalResults int      `json:"total_results"`
}

func SearchTVShows(ctx context.Context, query string) ([]Series, error) {
	res, err := request(ctx, "/search/tv?query="+url.QueryEscape(query))

	if err != nil {
		return nil, fmt.Errorf("error while trying to search for tv shows: %w", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error while trying to search for tv shows: unexpected status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("error while trying to search for tv shows: %w", err)
	}

	var s response

	if err := json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("error while trying to search for tv shows: %w", err)
	}

	return s.Results, nil
}
