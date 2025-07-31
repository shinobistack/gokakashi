package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	client    *http.Client
	baseURL   *url.URL
	userAgent *string

	common service
	Agent  *AgentService
	Scan   *ScanService
}

type service struct {
	client *Client
}

func New(baseURL string, client *http.Client) (*Client, error) {
	var err error
	c := &Client{}

	if client == nil {
		client = &http.Client{}
	}
	c.baseURL, err = url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	c.client = client

	c.common.client = c
	c.Agent = (*AgentService)(&c.common)
	c.Scan = (*ScanService)(&c.common)

	return c, nil
}

type RequestOption func(*http.Request)

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body any, opts ...RequestOption) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.userAgent != nil {
		req.Header.Set("User-Agent", *c.userAgent)
	}

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v any) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decErr := json.NewDecoder(resp.Body).Decode(v)
	if decErr == io.EOF {
		decErr = nil // ignore EOF errors caused by empty response body
	}
	if decErr != nil {
		err = decErr
	}
	return err
}
