package themoviedb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/url"

	"github.com/svenliebig/seq"
)

type response struct {
	Page         int      `json:"page"`
	Results      []Series `json:"results"`
	TotalPages   int      `json:"total_pages"`
	TotalResults int      `json:"total_results"`
}

type SearchTVQuery struct {
	page                int
	year                int
	first_air_date_year int
	include_adult       bool
	language            string
	Query               string
}

func (q SearchTVQuery) String() (r string) {
	r = "?"

	if q.page != 0 {
		r += fmt.Sprintf("page=%d&", q.page)
	}

	if q.year != 0 {
		r += fmt.Sprintf("year=%d&", q.year)
	}

	if q.first_air_date_year != 0 {
		r += fmt.Sprintf("first_air_date_year=%d&", q.first_air_date_year)
	}

	if q.include_adult {
		r += "include_adult=true&"
	}

	if q.language != "" {
		r += fmt.Sprintf("language=%s&", q.language)
	}

	if q.Query != "" {
		r += fmt.Sprintf("query=%s&", url.QueryEscape(q.Query))
	}

	return
}

func SearchTVShows(ctx context.Context, query SearchTVQuery) (s response, err error) {
	res, err := request(ctx, "/search/tv"+query.String())

	if err != nil {
		err = fmt.Errorf("error while trying to search for tv shows: %w", err)
		return
	}

	if res.StatusCode != 200 {
		err = fmt.Errorf("error while trying to search for tv shows: unexpected status code %d", res.StatusCode)
		return
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		err = fmt.Errorf("error while trying to search for tv shows: %w", err)
		return
	}

	if err := json.Unmarshal(body, &s); err != nil {
		err = fmt.Errorf("error while trying to search for tv shows: %w", err)
	}

	return
}

type searchTVShowsSeq struct {
	page, totalPages, totalResults int
	ctx                            context.Context
	query                          SearchTVQuery
}

func SearchTVShowsSeq(ctx context.Context, query SearchTVQuery) seq.Seq[Series] {
	return searchTVShowsSeq{ctx: ctx, query: query}
}

func (p searchTVShowsSeq) Iterator() iter.Seq2[Series, error] {
	return func(yield func(Series, error) bool) {
		if p.page == 0 {
			p.page = 1
		}

		for {
			p.query.page = p.page
			res, err := SearchTVShows(context.Background(), p.query)

			if err != nil {
				yield(Series{}, err)
				return
			}

			for _, s := range res.Results {
				if !yield(s, nil) {
					return
				}
			}

			p.page++

			if p.page > res.TotalPages {
				return
			}
		}
	}
}
