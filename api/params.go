package api

import (
	"fmt"
	"net/url"
)

type Params struct {
	Key   string
	Value string
}

func ExtractQueryParams(rawURL string) ([]Params, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	query := parsedURL.Query()
	params := make([]Params, 0, len(query))

	for key, values := range query {
		if len(values) > 0 {
			params = append(params, Params{Key: key, Value: values[0]})
		}
	}

	return params, nil
}

func BuildQueryString(params []Params) string {
	if len(params) == 0 {
		return ""
	}

	values := url.Values{}
	for _, p := range params {
		values.Set(p.Key, p.Value)
	}

	return "?" + values.Encode()
}
