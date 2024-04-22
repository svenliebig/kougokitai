package rest

import (
	"context"
	"fmt"
	"net/http"
)

var (
	_ error = &ErrNewRequest{}
)

type ErrNewRequest struct {
	url    string
	method string
}

func (err *ErrNewRequest) Error() string {
	return fmt.Sprintf("error while trying to create a new %q request with the url %q", err.method, err.url)
}

type ErrGet struct {
	url     string
	headers http.Header
}

func (err *ErrGet) Error() string {
	return fmt.Sprintf("error while trying to perform a 'GET' request to %q", err.url)
}

func Get(ctx context.Context, url string, o ...Option) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("%w,  %s", &ErrNewRequest{url, "GET"}, err)
	}

	for _, option := range o {
		option.apply(req)
	}

	response, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("%w,  %s", &ErrGet{url: url, headers: req.Header}, err)
	}

	return response, nil
}
