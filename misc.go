package usgs

import (
	"net/http"
)

type Option interface {
	apply(*Client)
}

type withHTTPClient struct {
	client *http.Client
}

var _ Option = (*withHTTPClient)(nil)

func (whc withHTTPClient) apply(c *Client) {
	c.httpClient = whc.client
}

func WithHTTPClient(c *http.Client) Option {
	return withHTTPClient{client: c}
}
