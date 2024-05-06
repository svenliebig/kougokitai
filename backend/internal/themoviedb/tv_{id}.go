package themoviedb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/svenliebig/query"
)

type TvIdResponse struct {
	Series
}

// episode_groups

type TvIdQuery struct {
	append_to_response string `query:"append_to_response"`
}

func (q TvIdQuery) String() (r string) {
	return "?" + query.Stringify(q, query.SkipEmpty{})
}

func (c *client) TvId(ctx context.Context, seriesId int, query TvIdQuery) (s TvIdResponse, err error) {
	res, err := c.request(ctx, "/tv/"+strconv.Itoa(seriesId))

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
