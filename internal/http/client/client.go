package client

import (
	"fmt"
	"net/http"
)

type Client struct {
	client  *http.Client
	token   string
	headers map[string]string
}

func New(opts ...Option) *Client {
	c := &Client{
		client: http.DefaultClient,
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	return c.client.Do(req)
}
