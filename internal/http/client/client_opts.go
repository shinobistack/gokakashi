package client

import "net/http"

type Option func(*Client)

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

func WithToken(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(c *Client) {
		c.headers = headers
	}
}
