package client

import (
	"fmt"
	"net/http"
	"os"
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

	cfClientID := os.Getenv("CF_ACCESS_CLIENT_ID")
	cfClientSecret := os.Getenv("CF_ACCESS_CLIENT_SECRET")
	if cfClientID != "" && cfClientSecret != "" {
		req.Header.Set("CF-Access-Client-Id", cfClientID)
		req.Header.Set("CF-Access-Client-Secret", cfClientSecret)
	}

	return c.client.Do(req)
}
